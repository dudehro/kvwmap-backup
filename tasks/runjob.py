#!/bin/python3
#requires python3.5 or later

import subprocess
import sys
import time
from datetime import datetime
import job_funcs
import os

defFile = job_funcs.get_configFileAbsPath(sys.argv[1])
runJob = sys.argv[2]

#make Backupdirectory
job_funcs.mkDirs(job_funcs.get_Backupdir(defFile))

job = job_funcs.get_jobDefinition(defFile, runJob)

while job['name'] != "":
    try:
        job_funcs.writeLog(defFile, job['name'], starttime=datetime.now().strftime("%s") )
        output = subprocess.run(job['command'], stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, check=True)
        job_funcs.writeLog(defFile, job['name'], endtime=datetime.now().strftime("%s"), exitcode=output.returncode, stdout=output.stdout, stderr=output.stderr, args=output.args)
        if 'start-on-success' in job.keys():
            job = job_funcs.get_jobDefinition(defFile, job['start-on-success'])
        else:
            job['name'] = ''
    except subprocess.CalledProcessError as output:
        job_funcs.writeLog(defFile, runJob, endtime=datetime.now().strftime("%s"), exitcode=output.returncode, stdout=output.stdout, stderr=output.stderr, args=output.args)
        if 'start-on-error' in job.keys():
            job = job_funcs.get_jobDefinition(defFile, job['start-on-error'])
        else:
            job['name'] = ''
