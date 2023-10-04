# instagram-insights
Discover various insights from your Instagram data.

## Use-cases
- Find out which instagram users are not following back
- Export followers and following lists

## Prerequisites
Complete all steps from this section.

### Request to download your Instagram data
You need to get a copy of your Instagram data using the [Data Download](https://www.instagram.com/download/request) tool. This is a manual process that cannot 
be automated. Instagram only allows their mobile application and website to perform this. While requesting a download of
your data, select the following options:
- Types of information:
  - [x] Followers and following
- Format:
  - [x] JSON
- Date range:
  - [x] All time

### Prepare your Instagram data
When your Instagram data is ready to be downloaded, you will receive an email notification.  
After downloading the zip archive, containing your Instagram data, you will have two options, depending on your use-case:
- [ ] Option 1: [Run the application on your machine](#run-the-application-on-your-machine)
- [ ] Option 2: [Run the application using GitHub Actions workflow](#run-the-application-using-github-actions-workflow)

## Run the application using GitHub Actions workflow
- [ ] Upload the zip file, containing your Instagram data, to a cloud storage service. Afterward, generate a public link to
the archive. The most common cloud storage service is probably Google Drive. Following the steps below will make your 
Instagram data archive publicly accessible:
  - Right-click: `Share`
  - General access: `Anyone with the link`
  - Role: `Viewer`
  - Left-click: `Copy link`
  - Left-click: `Done`
- [ ] Set up your repository
  - Fork the [cecobask/instagram-insights](https://github.com/cecobask/instagram-insights) repository to your account
  - Enable the `insights` workflow ([help](https://docs.github.com/en/actions/using-workflows/disabling-and-enabling-a-workflow))
  - Create a new repository secret ([help](https://docs.github.com/en/actions/security-guides/using-secrets-in-github-actions#creating-secrets-for-a-repository)):
    - Name: `ARCHIVE_URL`
    - Secret: this must be equal to the value of your public archive url
- [ ] Run the `insights` workflow ([help](https://docs.github.com/en/actions/using-workflows/manually-running-a-workflow))
  - Shortly, check the console output of the workflow
  - If any users are not following back you will see who they are

## Run the application on your machine
**TBA**
