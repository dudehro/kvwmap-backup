#!/bin/python3
#requires python3.5 or later

import subprocess
import sys
import time
from datetime import datetime
import job_funcs
import os

# gets called with name of job to start

# logs execution details of job to job.json
    #start-, endtime, exitcode, stdout, stderr

# writes log to helperfuncs.get_backupdir/jobslog.json

defFile = job_funcs.get_configFileAbsPath(sys.argv[1])
jobName = sys.argv[2]

#get Job
jobDef = job_funcs.get_jobDefinition(defFile, jobName)

#make Backupdirectory
job_funcs.mkDirs(job_funcs.get_Backupdir(defFile))
job_funcs.writeLog(defFile, jobName, starttime=datetime.now().strftime("%s") )

#execute and log
try:
    output = subprocess.run(jobDef['command'], stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, check=True)
    job_funcs.writeLog(defFile, jobName, endtime=datetime.now().strftime("%s"), exitcode=output.returncode, stdout=output.stdout, stderr=output.stderr, args=output.args)
except subprocess.CalledProcessError as output:
    job_funcs.writeLog(defFile, jobName, endtime=datetime.now().strftime("%s"), exitcode=output.returncode, stdout=output.stdout, stderr=output.stderr, args=output.args)

