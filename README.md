UploadX (Cloud Run edition)
======

A simple **server** written in Go for screenshot uploading (originally written for ShareX, works with alot of screenshot programs). Utilises the POST functionality for very speedy uploading.

The [original version is here](https://github.com/tombowditch/UploadX), this version is modified to run on [Google Cloud Run](https://cloud.google.com/run/) with [Google Cloud Storage](https://cloud.google.com/storage/) as a backend.

Usage
========

* Install and authenticate to [gcloud cli utilities](https://cloud.google.com/sdk/gcloud/)
* Clone repo
* Run `gcloud builds submit --tag gcr.io/YOUR_PROJECT/uploadx` (this may take a few minutes, this builds the container and uploads it to gcr.io)
* Whilst that's running, go to [console.cloud.google.com/storage](https://console.cloud.google.com/storage) and create a bucket. Choose a multi-region location (us or eu, normally). Keep Storage Class as 'Standard'. Keep everything else as default.
* OPTIONAL: Create a lifecycle rule to automatically delete objects from that bucket after X days. I have my bucket auto-delete screenshots after 7 days.
* Once your build is done, go to [console.cloud.google.com/run](https://console.cloud.google.com/run) (make sure the correct project is selected) and hit 'Create Service', select the correct container and your region.
* Expand the advanced section at the bottom and fill in three environment variables: `UPLOAD_KEY` (used for authentication in the POST request), `BUCKET_NAME` (your Google Cloud Storage bucket name), `PROJECT_ID` (your Google Cloud project ID)
* You're deployed. Be sure to add a domain mapping to use your custom domain, vs the supplied `a.run.app` subdomain.

Client libraries
================

An example bash script for macOS is available in `clientexample/screeny.sh`.
It is possible to use UploadX with ShareX.
