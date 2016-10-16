
This isn't a readme file -- it's more like a personal diary!

Now that I found a successful path (even if not yet an optimal and correct path) to the
creation of a Kibana Visualization using time-based data from my hackybeat, I wanted to
jot down a diary of what I did so that I can pick things up from here next time.



kill kibana, logstash, hackybeat

leave elasticsearch running but delete all indexes to start fresh

do as shown in this file as of the specified commit:
    git show fa33320d4e:scripts/README

cd into etc/kibana/

invoke etc/kibana/load.sh (from commit 07f591366ea) to load the index-pattern
(Note: load.sh won't find any 'search', 'visualization', or 'dashboard' json, but it will
 succeed on 'index pattern logstash-*')

restart kibana, logstash, hackybeat

go to kibana. it should prompt for 'Configure an index pattern'. accept defaults.

kibana should then show a page with 'This page lists every field in the logstash-* index'

    it should (among other things) contain:

      author.Email        string (analyzed & indexed)
      author.Email.raw        string (indexed)

      author.Name        string (analyzed & indexed)
      author.Name.raw        string (indexed)

      author.Uri        string (analyzed & indexed)
      author.Uri.raw        string (indexed)

      categories        string (analyzed & indexed)
      categories.raw        string (indexed)

      title        string (analyzed & indexed)
      title.raw        string (indexed)


it is then possible to make a visualization by categories. it looks like this:

    http://screencast.com/t/2pL61nY6he

    (Except now it looks BETTER than the screencast, due to commit 2c1bab120496c617a81e71bc51f0c982773d8009)
