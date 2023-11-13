# instagram-insights
Discover various insights from your Instagram data.

[![quality](https://github.com/cecobask/instagram-insights/actions/workflows/quality.yaml/badge.svg)](https://github.com/cecobask/instagram-insights/actions/workflows/quality.yaml)
[![codecov](https://codecov.io/gh/cecobask/instagram-insights/graph/badge.svg)](https://codecov.io/gh/cecobask/instagram-insights)

## Use-cases
- Find out which users are not following you back
- Export followers and following user lists in various formats (table, json, yaml)
- Sort the results by various criteria in ascending or descending order
- Set sorting criteria and order direction of the results
- Limit the number of results to get a quick overview (e.g. top 10)

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
- [ ] Option 1: [Run the application using GitHub Actions workflow](#run-the-application-using-github-actions-workflow)
- [ ] Option 2: [Run the application on your machine](#run-the-application-on-your-machine)

## Run the application using GitHub Actions workflow
- [ ] Upload the zip file, containing your Instagram data, to a cloud storage service. Afterward, generate a public link to
the archive. The most common cloud storage service is probably Google Drive. Following the steps below will make your 
Instagram data archive publicly accessible:
  - Right-click your archive and open the `Share` menu
  - In the `General access` section use the dropdown menu and select `Anyone with the link`
  - Select role `Viewer`
  - Use the `Copy link` button to copy the public link to your archive
- [ ] Set up your repository
  - Fork the [cecobask/instagram-insights](https://github.com/cecobask/instagram-insights) repository to your account
  - Enable the `insights` workflow ([help](https://docs.github.com/en/actions/using-workflows/disabling-and-enabling-a-workflow))
  - Create a new repository secret ([help](https://docs.github.com/en/actions/security-guides/using-secrets-in-github-actions#creating-secrets-for-a-repository)):
    - Name: `ARCHIVE_URL`
    - Secret: this must be equal to the value of your public archive url (_previously copied_)
- [ ] Run the `insights` workflow ([help](https://docs.github.com/en/actions/using-workflows/manually-running-a-workflow))
  - Shortly, check the console output of the workflow for results

## Run the application on your machine
- [ ] Download and install [Git](https://git-scm.com/downloads) + [Go](https://go.dev/doc/install)
- [ ] Clone the repository
- [ ] Build the application: `make build`
- [ ] Add the application to your path
- [ ] Load your Instagram data from a local zip file or cloud storage: `instagram information load <source>`
- [ ] Discover the available commands or browse through the [documentation](docs/instagram.md): `instagram --help`
