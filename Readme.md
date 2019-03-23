Waterloo Weather
================

Computes the daily low/high for a set of weather input, and plots a
recent-biased average low/high for each day of the year along with the
range of temperatures at which you should consider switching your
tires.

The data comes from http://weather.uwaterloo.ca

The source code is provided under the MIT license.



Instructions
------------

    $ make && ./ComputeLows data/*.csv && open weather.png
    go build ./...
