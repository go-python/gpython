# Copyright 2022 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

print('''
Welcome to a gpython embedded example, 
    where your wildest Go-based python dreams come true!''')

# This is a model for a public/user-side script that you or users would maintain,
# offering an open canvas to drive app behavior, customization, or anything you can dream up.
#
# Modules you offer for consumption can also serve to document such things.
from mylib import *

springBreak = Vacation("Spring Break", Stop("Miami, Florida", 7), Stop("Mallorca, Spain", 3))
springBreak.AddStops(Stop("Ibiza, Spain", 14), Stop("Monaco", 12))
springBreak.PrintItinerary()

print("\nI bet %s will be the best!\n" % springBreak.GetStop(4).Get()[0])
