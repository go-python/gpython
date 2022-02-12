# Copyright 2022 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import mylib_go as _go

PY_VERSION  = _go.PY_VERSION


print('''
==========================================================
        %s
==========================================================
''' % (PY_VERSION, ))


def Stop(location, num_nights = 2):
    return _go.VacationStop_new(location, num_nights)

    
class Vacation:

    def __init__(self, tripName, *stops):
        self._v, self._libVers = _go.Vacation_new()
        self.tripName = tripName
        self.AddStops(*stops)
        
    def __str__(self):
        return "%s, %d stop(s)" % (self.tripName, self.NumStops())
        
    def NumStops(self):
        return self._v.num_stops()  
        
    def GetStop(self, stop_num):
        return self._v.get_stop(stop_num)

    def AddStops(self, *stops):
        self._v.add_stops(stops)        
        
    def PrintItinerary(self):
        print(self.tripName, "itinerary:")
        i = 1
        while 1:
        
            try:
                stop = self.GetStop(i)
            except IndexError:
                break
                
            print("    Stop %d:  %s" % (i,  str(stop)))
            i += 1
            
        print("###  Made with %s " % self._libVers)
