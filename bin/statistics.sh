#!/bin/sh
cd "$(dirname "${0}")/.."

fmt="%-10s: "

printf "${fmt}" "stations"
{
  for d in htdocs/stations/*
  do
    [ -d "${d}" ] && echo "1"
  done
} | wc -l

printf "${fmt}" "first"
{
  for d in htdocs/stations/*/????/??
  do
    ls "${d}"/??/????\ *.xml | cut -d / -f 4-
  done
} | sort | head -n 1

printf "${fmt}" "last"
{
  for d in htdocs/stations/*/????/??
  do
    ls "${d}"/??/????\ *.xml | cut -d / -f 4-
  done
} | sort | tail -n 1

printf "${fmt}" "days"
{
  for d in htdocs/stations/*/????/??
  do
    ls -d "${d}"/??
  done
} | wc -l

printf "${fmt}" "broadcasts"
{
  for d in htdocs/stations/*/????/??
  do
    ls "${d}"/??/????\ *.xml
  done
} | wc -l

printf "${fmt}" "podcasts"
{
  ls -d htdocs/podcasts/*
} | wc -l

printf "${fmt}" "recordings"
{
  for d in htdocs/enclosures/*/????/??
  do
    ls "${d}"/??/????\ *.* 2>/dev/null
  done
} | wc -l

echo "Disk usage..."
echo " index.xml (day summary, gzip) : $(du -hsc htdocs/stations/*/????/??/??/index.xml.gz | fgrep total | cut -f 1)"
# echo " non-index.xml (single broadcasts): $(ls -f htdocs/stations/*/????/??/??/index.xml | cut -d / -f 1-6 | xargs du -ksc --exclude index.xml --exclude='*.json' | fgrep total | cut -f 1) KiB"
echo " broadcasts                    : $(du -hsc htdocs/stations | tail -n 1 | cut -f 1)"
echo " broadcasts, no index          : $(du -hsc --exclude 'index.xml.gz'                    htdocs/stations/*/???? | tail -n 1 | cut -f 1)"
echo " broadcasts, no index, no json : $(du -hsc --exclude 'index.xml.gz' --exclude='*.json' htdocs/stations/*/???? | tail -n 1 | cut -f 1)"
