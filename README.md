## Fun with Elastic Beats

Built on top of [https://github.com/elastic/beats](https://github.com/elastic/beats), with inspiration from [the community beats](https://www.elastic.co/guide/en/beats/libbeat/master/community-beats.html).

### Goals

- Implement at least 2 new beats...
- ...Then iterate and refine the code to *extract as much duplication as humanly possible* (both literal and conceptual duplication)

### Results (thus far)

- Two significant areas of duplication have been identified.
- Proof-of-concept reusable components are implemented in the hackybeat repo.  Such reusable components would address duplicate work normally reinvented in each new beat.

### Details

The 2 beats in this repo are the sample `rss` beat and the sample `gitlog` beat.

As I implemented these 2 beats, I noticed that they each needed some kind of `ticker` loop, and they each needed some kind of persistence/caching for "de-dupe" (deduplication) purposes.

Those two needs are not unique to my sample beats. [Every community beat](https://www.elastic.co/guide/en/beats/libbeat/master/community-beats.html) implements its own `ticker` or some form of looping.

Several community beats would benefit from a built-in de-dupe service.  Beats that need de-duping are those that periodically open and read from the same file or feed, as opposed to beats that always take one discrete point-in-time measurement during their hearbeat.

### Proof of Concept

This repository contains a reusable `util/poller` component and a reusable `util/deduper` component.  These components transparently benefit the `rss` sample beat and the `gitlog` sample beat *without the beats needing to know anything about it.*

![architecture diagram](https://raw.githubusercontent.com/pestophagous/hackybeat/master/doc/arch_dep.png)

Architectural dependency diagram powered in part by: [https://github.com/kisielk/godepgraph](https://github.com/kisielk/godepgraph)

----------

### Questions?

You know how to reach me :)

If you don't know how, then here's one way: http://bit.ly/25XJBhy 
