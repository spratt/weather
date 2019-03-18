Waterloo Weather
================

This is a quick hacky app that computes the yearly low, yearly high,
and the first day each year that the daily high stays above 7 degrees
Celsius in Waterloo, Ontario.  The last piece of information may be
useful for anyone trying to decide when to change their winter tires
to summer tires.

The data comes from http://weather.uwaterloo.ca

The source code is provided under the MIT license.



Instructions
------------

    $ make && ./ComputeLows data/*.csv
    go build ./...
    data/1999_weather_station_data.csv
    data/2000_weather_station_data.csv
    data/2001_weather_station_data.csv
    data/2002_weather_station_data.csv
    data/2004_weather_station_data.csv
    data/2006_weather_station_data.csv
    data/2007_weather_station_data.csv
    data/2008_weather_station_data.csv
    data/2009_weather_station_data.csv
    data/2010_weather_station_data.csv
    data/2011_weather_station_data.csv
    data/2012_weather_station_data.csv
    data/2013_weather_station_data.csv
    1999,-25.370000,32.800000,[1999-04-18 00:00:00 +0000 UTC]
    2000,-24.620000,30.370000,[2000-04-19 00:00:00 +0000 UTC]
    2001,-24.460000,35.350000,[2001-04-20 00:00:00 +0000 UTC]
    2002,-22.730000,32.990000,[2002-05-22 00:00:00 +0000 UTC]
    2004,-28.670000,29.310000,[2004-05-10 00:00:00 +0000 UTC]
    2006,-19.700000,33.650000,[2006-04-10 00:00:00 +0000 UTC]
    2007,-26.750000,33.130000,[2007-04-18 00:00:00 +0000 UTC]
    2008,-23.380000,30.400000,[2008-04-15 00:00:00 +0000 UTC]
    2009,-28.750000,31.290000,[2009-04-24 00:00:00 +0000 UTC]
    2010,-21.490000,33.121000,[2010-05-28 00:00:00 +0000 UTC]
    2011,-28.841660,38.037510,[2011-05-01 00:00:00 +0000 UTC]
    2012,-18.310699,34.023830,[2012-05-12 00:00:00 +0000 UTC]
    2013,-24.561251,34.718170,[2013-05-14 00:00:00 +0000 UTC]
