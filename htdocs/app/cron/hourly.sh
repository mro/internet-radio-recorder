#!/bin/dash
#
# Copyright (c) 2013-2016 Marcus Rohrmoser, http://purl.mro.name/recorder
#
# Permission is hereby granted, free of charge, to any person obtaining a copy of this software and
# associated documentation files (the "Software"), to deal in the Software without restriction,
# including without limitation the rights to use, copy, modify, merge, publish, distribute,
# sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all copies or
# substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT
# NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
# NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES
# OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
# CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
#
# MIT License http://opensource.org/licenses/MIT

cd "$(dirname "$0")"/../..
me="$(basename "$0")"

echo "Start  $(/bin/date +'%FT%T')" 2>> log/"${me}".stderr.log | tee -a log/"${me}".stdout.log

cmd=../bin/scrape-$(uname -s)-$(uname -m)-0.2.1

[ -x "${cmd}" ] || { echo "Executable ${cmd} not found." 1>&2 && exit 1 ; }

${cmd} 2>> log/"${me}".stderr.log \
| tee log/"${me}".stdout.dat \
| app/broadcast-render.lua --stdin 2>> log/"${me}".stderr.log \
1>> log/"${me}".stdout.log

nice app/calendar.lua stations/* podcasts/* 1>> log/"${me}".stdout.log 2>> log/"${me}".stderr.log

echo "Finish $(/bin/date +'%FT%T')" 2>> log/"${me}".stderr.log | tee -a log/"${me}".stdout.log
