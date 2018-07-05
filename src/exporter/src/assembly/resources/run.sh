#!/bin/sh
DIRNAME=`dirname $0`
RUNHOME=`cd $DIRNAME/; pwd`


echo @RUNHOME@ $RUNHOME
cd $RUNHOME


echo "\n\n### Starting pgExporter-go";
$RUNHOME/exporter




