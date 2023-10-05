#!/bin/python3
#requires python3.5 or later

import subprocess
import sys
import time
from datetime import datetime
import job_funcs
import os

def log2Stdout(str):
    print(f"{datetime.now().strftime('%F %r')}:  {str}")

def startJob(defFile, runJob, workdir):
    failedJobs = 0
    jobQ = list()
    jobQ.append( job_funcs.get_jobDefinition(defFile, runJob) )

    while len(jobQ) > 0:
        job = jobQ.pop()
        log2Stdout(f"running job: {job['name']}")
        try:
            job_funcs.writeLog(workdir, job['name'], starttime=datetime.now().strftime("%s") )
            if 'next-job' in job.keys():
                jobQ.append( job_funcs.get_jobDefinition(defFile, job['next-job']) )
            env=os.environ.update({'WORKDIR': workdir})
            output = subprocess.run(job['command'], env=env, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, check=True)
            job_funcs.writeLog(workdir, job['name'], endtime=datetime.now().strftime("%s"), exitcode=output.returncode, stdout=output.stdout, stderr=output.stderr, args=output.args)
            if 'start-job-on-success' in job.keys():
                jobQ.append( job_funcs.get_jobDefinition(defFile, job['start-job-on-success']) )
        except subprocess.CalledProcessError as output:
            failedJobs+=1
            job_funcs.writeLog(workdir, job['name'], endtime=datetime.now().strftime("%s"), exitcode=output.returncode, stdout=output.stdout, stderr=output.stderr, args=output.args)
            log2Stdout(f"job failed with error: {output.stderr}")
            if 'start-job-on-error' in job.keys():
                jobQ.append( job_funcs.get_jobDefinition(defFile, job['start-job-on-error']) )
            if 'exit-queue-on-error' in job.keys():
                jobQ = list()

    return failedJobs

defFile = job_funcs.get_configFileAbsPath(sys.argv[1])
runJob = sys.argv[2]

workdir = job_funcs.get_Workdir(defFile)
job_funcs.mkDirs(workdir)
logFile = os.path.join(workdir, 'joblog.json')

try:
    log2Stdout(f"definition: {defFile}")
    log2Stdout(f"logfile: {logFile}")
    failedJobs = startJob(defFile, runJob, workdir)
    log2Stdout(f"{failedJobs} jobs failed")
    sys.exit(failedJobs)
except Exception as e:
    print(f"Abbruch mit Fehler: {str(e)}")
    sys.exit(1)
