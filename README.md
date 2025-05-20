# Adventures in permutations

**The is a side quest**

I was exploring Genetic Algorithms when I got to the Travelling Salesman Problem. This is a classic problem to solve with GAs. After trying some of the smaller problems, 5 to 7 cities, I wanted to test out larger problems. Maybe 20 cities

The problem was without knowing what the shortest path was I had no idea how well my GA solution was and finding non trivial examples with known solutions was proving fruitless. So I though "I'll just brute force it!"

This does work for a small number of cities but with each additional city it got longer and longer. Brute force needs to examine factorial N routes. For N=16 that is 20,922,789,888,000 (twenty trillion, nine hundred and twenty-two billion, seven hundred and eighty-nine million, eight hundred and eighty-eight thousand) and can take days just to calculate all the permutations let alone the distance travelled between the cites

There are things that can be done to speed things up. Given that the route is circular 1 -> 2 -> 3 -> 1 (for N=3) there is no need to examine 2 -> 3 -> 1 -> 2 or 3 -> 1 -> 2 -> 3 as they are all the same route. So you pick a root city, say 1, and just don't check the rest. Instead of checking N! routes we are now checking (N-1)! routes. A speed up for checking but you still need to *generate* the N! routes. Below is a table of how long it took me to simply generate the number of permutations...

|Element|Permutations|Time|Notes|
|---|--:|--:|---|
|5|120|0.00001658s||
|6|720|0.00000487s||
|7|5,040|0.00001596s||
|8|40,320|0.00010342s||
|9|362,880|0.00083308s||
|10|3,628,800|0.00786821s||
|11|39,916,800|0.06768887s||
|12|479,001,600|0.76926550s||
|13|6,227,020,800|10s||
|14|87,178,291,200|2m 22s||
|15|1,307,674,368,000|35m 27s||
|16|20,922,789,888,000|9h 25m||
|17|355,687,428,096,000|6d 16h|Estimated|
|18|6,402,373,705,728,000|120d 1h|Estimated|
|19|121,645,100,408,832,000|6y|Estimated|
|20|2,432,902,008,176,640,000|125y|Estimated|

The problem with calculating permutations is that you cannot hand off the first 1,000,000 calculations to one cpu and the second 1,000,000 to another. The second cpu does not know where to start, what is the 1,000,001st permutation without calculating the preceding 1,000,000?

Is there a way to distribute the calculation of permutations to use as much hardware as is available?

Maybe...

This will sound a little odd. But to calculate all the permutations for 6 cities couldn't I just calculate for cities 1 to 3 and then for 4 to 6 separately and then _somehow_ merge them?

Well oddly the answer is yes

So how many ways are there to mangle two list? In this case we ant to find all the permutations for 6 elements by mangling two permutations of 3 elements. Here is how we can combine the left with the right

	LLLRRR
	LLRLRR
	LRLLRR
	RLLLRR
	LLRRLR
	LRLRLR
	RLLRLR
	LRRLLR
	RLRLLR
	RRLLLR
	LLRRRL
	LRLRRL
	RLLRRL
	LRRLRL
	RLRLRL
	RRLLRL
	LRRRLL
	RLRRLL
	RRLRLL
	RRRLLL

So if we have `3,1,2` from the left (first) list and `4,6,5` from the right (second) list then there are 20 ways they can be combined. For example we have `3,1,2` and `4,6,5` and `LRRLLR`. We replace the `L` with values from the left list and the `R`s with values from the right list and we get `3,4,6,1,2,5`

So I wrote a quick test program and to generate all the permutation of 12 elements by merging two sets of 6 permutations worked and only took 4.45637037 seconds! Which is somewhat slower than the 0.76926550 seconds that it took before but it worked and despite being slower has a significant advantage in that we can now split the work out to as many cpus / cores we can get our hands on. We can also stop processing and restart as we wish



