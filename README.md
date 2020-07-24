# Iterum CLI

---

If you haven't done so already, check out the [main repository](https://github.com/iterum-provenance/iterum) with introductory material.

---

The CLI for the iterum tool is the main user interface. After deploying the relevant Iterum infrastructure the CLI can be used to interact with those components in order to version data sets, create pipelines and evaluate their results. The CLI is divided into two main sub-components. On the one hand there is the data version control, and on the other hand pipeline interaction. This split is kept strict such that at a later stage their functionalities can be split easily across multiple interfaces. The CLI divides the code into 3 main package: `cmd` for CLI commands structured with cobra, `pipeline` for pipeline related functionality and `idv` for data versioning related functionality

---

## Quick start

### Installing the CLI
* `make build` builds the go application in `./build/iterum`
* `make link` creates a symlink to the created build folder making `iterum` accessible via the terminal
* `make clean` removes the symlink pointing to `/usr/bin/iterum` and the build folder `./build`

### Deploy your first pipeline

---

## Iterum data versioning (IDV)
The Iterum data versioning component deals with versioning data. It consumes API features exposed by the [Iterum Daemon](https://github.com/iterum-provenance/daemon), and further operates on local file systems. Its commands are defined in `cmd/data` and exposed in the CLI using `iterum data [subcommand]`. 

### Concepts
IDV works similar to git in a few ways. There is a local version that you make changes on (similar to staging changes) and these changes can be committed (similar to commit & pushing). The remotely stored version is accessible via the Daemon, so IDV interacts with the Daemon, which in turn ensures how data is stored in persistent storage. These 3 layers allow for a complete decoupling of the components and other instances of a Daemon to mount the same datasets at later times. 

#### idv-config.yaml
The core of each (data versioning) repository is the idv-config. This file describes how and where a dataset is located. Ideally users only need to share their `idv-config.yaml` in order for others to be able to collaborate on their data sets. The naming of the file is important. A simple example is shown below
```yaml
name:                           # Name of the dataset
    test-cats
description:                    # Short description of the dataset
    images of cats all the way
daemon:                         # URL at which to reach the Daemon
    http://localhost:3000/
backend:                        # Where is the data stored
    Local                       
# { Local, AmazonS3, GoogleCloud }, though only Local is implemented
credentials:                    # backend-specific credentials
  path:                         # Local storage only needs a path
    /localStorage/             
```
#### .idv folder
Each IDV repository is accompanied by a `.idv`-folder. This folder contains all management related data for a repository. Its version tree, commit files, branches, etc. There are three main symlinks that get swapped around `BRANCH` (the current branch file), `LOCAL` (commit only known locally tracking staged changes) and `HEAD` (current remotely known HEAD of the tree).

#### How data is versioned
Iterum employs a crude data versioning component and assumes a user knows when he or she is adding something. Data files can be added from any accessible file storage on the user's machine. If files are new they are added, if files clash in naming with a known file within the data set it is considered an update of that file. Files are not sliced nor deleted from the remote storage, but stored based on the commit that introduces them. For details on the actual storage this repository refers to the [Daemon](https://github.com/iterum-provenance/daemon). Updates to files always completely replace older versions, so no diff-ing based on files. This makes individual updates to CSV files very inefficient.

### Usage
For setting up a new repository pick a folder and then:
1. Run `iterum data init`, this generates a simple idv-config.yaml.
2. Fill out the details of the config file.
3. Run `iterum data setup` to make the repository initialize and let the Daemon know of the dataset.

Now that the repo is set up the following commands can be used to interact with it
* Use `iterum data {add, remove, unstage}` to stage changes and files
* Use `iterum data {ls, status, inspect, download}` to analyze and interact with the dataset and its staged changes
* Use `iterum data commit "[tag]" "[description]"` to push changes to the daemon and make them persistent
* Use `iterum data {branch, checkout}` to switch branches and check out other commits 

---

## Iterum pipeline management (IPM)
The Iterum pipeline management component deals with pipeline deployment, inspection and results analysis. It consumes API features exposed mainly by the [Iterum Manager](https://github.com/iterum-provenance/manager), but some also by the [Iterum Daemon](https://github.com/iterum-provenance/daemon) due to its access to pipeline results. Its commands are defined in `cmd/pipeline` and are exposed in the CLI using `iterum pipeline [subcommand]`. 

### Concepts

#### Pipeline Run Configuration
Pipelines are defined using JSON structures that describe a DAG of transformations, fragmenter, combiner and data set version.
A simple example of which is shown below. The Manager that deploys this needs to have access to a Daemon that knows this data set as well.
Multiple publishers and consumers to the same channels are allowed, however note that messages are not replicated for each consumer (as of yet).
These pipeline run configurations can be submitted using the CLI as shown below. 
```json
{
  "name": "dummy-pipeline",
  "input_dataset": "dummy-data",
  "input_dataset_commit_hash": "SsWkbrzAF8PzXysRmOcjpnQrCcPHiO6B",
  "fragmenter": {
    "image": "localhost:32000/dummy-fragmenter:latest",
    "output_channel": "fragmenter_output"
  },
  "steps": [
      {
        "name": "new-step1", 
        "image": "localhost:32000/dummy-transformation-step:latest", 
        "input_channel": "fragmenter_output", 
        "output_channel": "step1_output"
      },
      {
        "name": "new-step2", 
        "image": "localhost:32000/dummy-transformation-step:latest", 
        "input_channel": "step1_output", 
        "output_channel": "step2_output"
      }
  ],
  "combiner": {
    "input_channel": "step2_output"
  }
}

```

#### Achieving DAG structures
Due to its early alpha stages, Iterum does not cleanly expose the option to defines complex DAGs in the JSON specification. However it does provide handles for doing so, by making transformation steps post fragment descriptions with `metadata.output_channel`, other output channels than the default can be defined. In [Pyterum](https://github.com/iterum-provenance/pyterum) this can be done using the [metadata structure](https://github.com/iterum-provenance/pyterum/blob/master/pyterum/metadata.py) in a fragment description. This functionality is not yet in a completed stage and should be on the next feature list.

### Usage
* Use `iterum pipeline submit [specification.json]` to to deploy and run a pipeline
* Use `iterum pipeline {delete, purge}` to remove kubernetes resources of pipelines or to completely purge the pipeline and its results
* Use `iterum pipeline {status, history}` to view the progress of current and past pipelines
* Use `iterum pipeline {download, results lineage}` to retrieve results of the pipeline. Either output data or lineage on the fragment level 

---

## Code documentation
The documentation of code is left to code files themselves. They have been set up to work with Godoc, which can be achieved by running `godoc` and navigating to `http://localhost:6060/pkg/github.com/iterum-provenance/cli/`.
