## Requirements

5M  questions
20M answers 
10M people access this site
50M page views per month, 19/second
Expecting to double in size within one year

Say on average, the content of a post is 2KB. Then:
10GB of questions (plus metadata)
40GB of answers   (plus metadata)

...where metadata includes:
- author (64-bit int, FK)
- rating (32-bit int)
- creation time (64-bit int).
For answers, it also includes 64-bit q_id.

We need to support the following queries:
- Serve the questions and answers for question q_id.
- Serve the top N questions sorted by creation time.
- Serve the top N questions sorted by rating.

## One machine

To serve a Q&A-thread, 2KB + (A * 2KB) must be sent to the client,
where it's 2KB average for the question, and A represents the number of answers for that question.

_Is it possible?_
Seems likely. We can store relations between questions and answers using a normal SQL database, as well as being able to sort on columns like rating and creation time using relational indexes.

_Is it feasible?_
At this scale, it probably is. Accepting HTTP requests and going to disk ~19 times per second, with storage of ~100GB including database indexes, is feasible. This could all be done with one webserver, and one Postgres instance with sufficient storage (256 GB, for example).
However, these machines would be single points of failure, and disk failures would mean loss of data since there is no backup. Additionally, as our memory usage increases (expecting to double this year alone), the only means of scaling up is by adding more memory to our Postgres machine.

_Can we do better?_
Yes, we can introduce overall fault tolerance, replication to avoid data loss, and a means by which to scale the system horizontally.

## Add replication

The first thing we want to avoid is data loss in the event of disk failure. We can use Postgres' built-in replication to back up the database to multiple machines. One machine should be for cold storage, backed up periodically, and the rest can serve as hot read-only replicas. This way, read-availability is preserved in the event that the main postgres node goes down, and maybe we can implement switching the writes to the secondary node in those cases, too, to preserve write availability (about 4x less valuable than read-availability, but still something to keep in mind, not sure how much more difficult this is or if Postgres makes it easy).

Additionally, the webservers are horizontally scalable by nature, so we can deploy a single load balancer to send traffic to webservers. At 19 requests per second, we are nowhere near hitting the limits of a single HTTP server, but we don't want a single machine failure to take down the whole site. So we can start with 2 web servers behind a load balancer, and scale as needed depending on how flaky our machines are.

In total, we now need:
- A load balancer
- 2 web servers
- 2 Postgres servers (one for writes, the other a read-only replica)
- 1 machine for cold storage of DB data, backed up periodically

_Is it feasible?_
This architecture can easily handle our current scale, and if we double our scale by next year, the system will not have to change at all to accommodate it.

_Can we do better?_
If we do continue to double in scale year over year, eventually having complete read-replicas for all Postgres nodes will not lend itself well to horizontal scaling. As it stands, we are accounting for having 256GB of disk storage per Postgres node -- to continue scaling, we'd need to continue adding more and more storage to each node. It would be nice if the Postgres nodes didn't actually need 256GB of storage each.

## Sharding and Caching

To reduce the necessary storage requirements of each Postgres node, we can store subsets of the data on each node, and route reads/writes to the correct machine based on the sharding key. Sharding on question ID seems natural, since answers have question id FKs and can be kept on the same shard as their parent (which makes the common query of Q&A-thread as efficient as possible). Even though some questions will have way more answers than others, the random distribution of questions across shards should eventually have equally popular questions being stored across different shards. We could get unlucky with this, but there is room for error -- each answer is about 2KB of storage. Will there really be questions with 100,000 answers (which would be a skew of (100,000*2KB) = 200MB for that question)?

_Is it feasible?_
Yes. You can do a simple hash-mod type of sharding, or use consistent hashing to allow for easier addition/subtraction of nodes. However, sharding makes some of the queries that we need to support potentially more complicated. Reminder of the queries we support:
- Serve the questions and answers for question q_id.
  * Still works fine with sharding.
- Serve the top N questions sorted by creation time.
  * Sharding on question ID makes this harder, since we will have to query every node for its top N most recent questions.
- Serve the top N questions sorted by rating.
  * Same as above.

_Can we do better?_
Separating the data necessary for querying ratings / creation time from the Postgres replicas might make the queries perform better.

## Caching recent requests

With the new architecture of sharding questions and answers based on question id, we'd like an efficient way to query for the top N most recent questions. To do this, whenever a new question is created, the webserver that handles the creation can send the newly created question id to a dedicated machine which simply stores question ids keyed by creation time. Since the "recent questions" query will likely not exceed ~100 questions, this RecencyCache can have quite modest storage available -- each entry will be 8 bytes for a questoin id, and 8 bytes for its creation timestamp. 16GB of storage can store a billion questions, easily sorting them by time by having transactional appends to a log (or even using a single-table SQL database, which might even be too heavyweight but potentially better than rolling our own atomic append-only log).

For queries about most recent questions, webservers can hit the RecencyCache for a list of the N most recent questions, and then use those ids to query the Postgres nodes for the full payload.

## Ratings

The final query to solve will be serving the top N questions sorted by rating. Currently, a webserver would need to query each Postgres node for its N highest rated questions, and once all results are in, compute the true N highest rated questions. Since all of these queries can run concurrently, its performance will be as slow as the slowest of the Postgres nodes. In truth, this is probably fine (even "failing open" if one of the Postgres nodes does not respond, and just serving the highest rated questions of the successful calls, since that is probably an acceptable UX, but gotta check with Product on that).

_Can we do better?_
Maybe. Similar to the RecencyCache, maybe there is an opportunity to store rating changes on a separate machine to avoid the scatter-gather approach required after sharding the data on question id.

Storing a question's (or answer's) rating as a mutable field in Postgres loses the rating history of all posts. So if in the future, for example, we wanted to show "trending questions", or some other query that relies on the history of a post's rating, we would not be able to do so. Instead, if we stored a log of rating change events, we could quickly determine the top ratings for questions (or answers) by keeping an in-memory hash table of post_id->rating, and if the system crashes, scanning through the event stream and re-tallying the totals.

For example, we can have this RatingEvent machine have two streams -- one for questions, one for answers. Each entry in the log will contain the id of the record, the time, and say one byte to determine the type of event (1 bit for determining upvote vs. downvote, 1 bit for question-or-answer, and the other bits reserved for possible extensions in the future). The machine would hold an in-memory tally of question ratings and answer ratings, updating as events come in, and if/when the machine goes down and reboots, it will scan the log and recompute the totals in memory. It can also checkpoint its calculations, move old event logs to cold storage, and incorporate new data into the last checkpoint as the starting state. To checkpoint the state of ratings for the entire system requires 12 bytes per post (we'll store a post's rating in a 32-bit int, imposing a limit of about 2 billion upvotes or downvotes per post). That means even with 1 billion posts having checkpointed state, the entire snapshot would fill 12GB, which would be scanned into memory on boot. 

_Is it feasible?_
If every one of our 50M monthly views we receive were to result in a rating change (which probably couldn't happen, since I think only authenticated users are allowed to vote), it would require about 1GB per month (8 bytes for the post id, 8 bytes for the time, and one byte to describe the type of event, say 20 bytes in total), or 12GB/year.

_Can we do better?_
Moving ratings to a separate event stream introduces a tradeoff:
- Keep the ratings in the regular Postgres nodes, and accept the scatter-gather nature of "top N rated questions/answers" queries, or
- Move the ratings to a separate event stream, and require that extra query every time we display information about a post (since Postgres doesn't store ratings anymore, but it's data that we'll need to show whenever we serve a Q&A thread).

The new performance of serving a basic Q&A thread is the query to the proper Postgres shard (located via the thread's question_id), and then another query to the RatingEvent stream with the associated question/answer ids, and finally serving that data to the client. The extra latency would be pretty minimal (in-memory lookups from the event stream, with persistent TCP connections between the webserver and the event stream), but the event stream introduces a single point of failure, and requires replication to improve availability, as well as adds logic and infrastructure to achieve performance that isn't required for this feature. Populating a dashboard with the top N highest rated queries is not as hot a path as serving an entire Q&A thread.

I'll need to check with product, but populating the dashboard by scatter-gather querying each Postgres node -- with a "fail-open" strategy on timeout/error -- seems acceptable, since it's only as slow as the slowest node, and the value of this query is probably not high enough to warrant the extra event stream. But we have a fallback if this option is not acceptable to our users (and the event stream can always be added later, by using the current rating as the starting state, and sending all new rating changes to the event stream).

## Overall architecture

- A load balancer (commodity hardware with limited memory which supports ~100 IOPS, 8 or 16GB of RAM is fine)
- 2 web servers (more can be added later to scale horizontally, similar specs to the load balancer, maybe higher RAM in case we want to start caching questions in the future -- say 32 or 64 GB RAM depending on price)
- 5 Postgres servers (arbitrarily picking 5 nodes for sharding, at the current storage of ~100GB each shard would be responsible for about 20GB, so let's say 64GB per node, with the option to scale up in the future if need be with more nodes and more shards)
- 1 machine for cold storage of DB data, backed up periodically (not sure how this normally works -- would each shard get its own cold storage backup, or would they all send to the same place? could just get multi-TB disks and take intermittent backups for all of the nodes, or get one per node and then each backup disk could be small... what's more common, restoring one shard, or restoring the entire DB? seems like backing up per shard makes more sense, since all shards sending data to one node requires coordination maybe? Or is there an established way to do that?).

Requests:

- User request to post a new question hits the load balancer, with the user's unique 8-byte id and the question text. The LB sends the request to an available webserver. The webserver generates a unique 8-byte question id, looks up the appropriate Postgres shard, and sends it the data. Postgres stores the question id, the content, the timestamp, the author id, and a rating of 0. The webserver also sends the question ID to the RecencyCache, which stores the id and timestamp of creation.

- Posting an answer to a question is basically the same flow, just with the additional field of answer_id (also passing the question_id as a FK). Answers are stored on the same shard as their question.

- Querying for a question and its associated answers hits the LB->webserver, sending the question_id desired. The webserver finds the appropriate Postgres shard, sends the query, and Postgres selects all answers which have a FK of the provided question_id, as well as the question itself. It sends this data back to the webserver, which passes it downstream to the client.

- Querying for the top 10 most recent questions hits the LB->webserver, the webserver hits the RecencyCache which takes its 10 most recent question_ids and sends them to the webserver. The webserver then queries the appropriate Postgres nodes, but only for the questions (the answers aren't necessary).

- Querying for the top 10 highest rated questions hits the LB->webserver. The webserver sends a request to every Postgres shard asking for its top 10 highest rated questions, with a timeout set to the max SLA for this query (minus a little buffer, which we can determine with Product). It then compares the result sets to determine the 10 globally highest rated questions, and serves those to the client.
