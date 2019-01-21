---
title: Werf Frequently Asked Questions
permalink: faq.html
layout: page
---

## General

[Q: Can I use different werf/docker version for build and for deploy?](#general-1){:id="general-1"}

In some case, you can and it will work, but please try to avoid this and use latest werf version.


## Werf config

[Q: How to convert **COPY . /var/app** instruction from Dockerfile to werf.yaml?](#config-1){:id="config-1"}

To add files from local git repository you can use the following:

```yaml
git:
- add: /
  to: /var/app
```


[Q: How to specify stageDependency to all files in all subdirectories?](#config-2){:id="config-2"}

You can use `**/*` mask.


[Q: How can I use environment variable in config?](#config-3){:id="config-3"}

Use [sprig env function](http://masterminds.github.io/sprig/os.html), `{% raw %}{{ env "ENV_NAME" }}{% endraw %}`, for werf.yaml:

{% raw %}
```yaml
{{ $_ := env "SPECIFIC_ENV_HERE" | set . "GitBranch" }}

dimg: ~
from: alpine
git:
- url: https://github.com/company/project1.git
  branch: {{ .GitBranch }}
  add: /
  to: /app/project1
- url: https://github.com/company/project2.git
  branch: {{ .GitBranch }}
  add: /
  to: /app/project2
```
{% endraw %}

[Q: Can I set an environment variable to use it during build?](#config-4){:id="config-4"}

We recommend to build an image which building instructions depend on your code but not on an environment in build time. In other words, you better build one image, which you can run in any environment, and this image has to change its logic when it runs rely on environment variables. If building stage will depend on such black box like changing environments you can get an unexpected behavior of werf builder and unexpected results.

Environment variables which have been set in `docker` config section will be added by a builder on the last dimg stage, `docker_instructions`, and will not be accessible on other build stages.

Also, you can use `WERF_ANSIBLE_ARGS` env when you use ansible builder. E.g. you can `export WERF_ANSIBLE_ARGS=-vvv` and get verbose ansible output.


[Q: What functions can I use in werf.yaml?](#config-5){:id="config-5"}

* [Go template functions](https://golang.org/pkg/text/template/#hdr-Functions).

* [Sprig functions](http://masterminds.github.io/sprig/).

* `include` function with `define` for reusing configs:

{% raw %}
```yaml
dimg: app1
from: alpine
ansible:
  beforeInstall:
  {{- include "(component) ruby" . }}
---
dimg: app2
from: alpine
ansible:
  beforeInstall:
  {{- include "(component) ruby" . }}

{{- define "(component) ruby" }}
  - command: gpg --keyserver hkp://keys.gnupg.net --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3
  - get_url:
      url: https://raw.githubusercontent.com/rvm/rvm/master/binscripts/rvm-installer
      dest: /tmp/rvm-installer
  - name: "Install rvm"
    command: bash -e /tmp/rvm-installer
  - name: "Install ruby 2.3.4"
    raw: bash -lec {{`{{ item | quote }}`}}
    with_items:
    - rvm install {{ .RubyVersion }}
    - rvm use --default 2.3.4
    - gem install bundler --no-ri --no-rdoc
    - rvm cleanup all
{{- end }}
```
{% endraw %}

* `.Files.Get` function for getting project file content:

{% raw %}
```yaml
dimg: app
from: alpine
ansible:
  setup:
  - name: "Setup /etc/nginx/nginx.conf"
    copy:
      content: |
{{ .Files.Get ".configs/nginx.conf" | indent 8 }}
      dest: /etc/nginx/nginx.conf
```
{% endraw %}


## Building

[Q: How to specify ssh keys?](#building-1){:id="building-1"}

Use `--ssh-key=path-to-id-rsa` option. E.g. `werf dimg build --ssh-key=path-to-id-rsa`.


[Q: I've added files from a git repository, but I can't access it](#building-2){:id="building-2"}

You can't access files on stage `before_install` because werf adds sources on stage `git_archive`, and therefore you can access files on any stage after (e.g `install`, `before_setup`, `setup`).


[Q: How can I build specific dimgs?](#building-3){:id="building-3"}

You can specify `DIMG`, dimg name, for most commands:
```bash
werf dimg build [options] [DIMG ...]
werf dimg bp [options] [DIMG ...] REPO
werf dimg push [options] [DIMG ...] REPO
werf dimg spush [options] [DIMG] REPO
werf dimg tag [options] [DIMG ...] REPO
werf dimg run [options] [DIMG] [DOCKER ARGS]
werf dimg stage image [options] [DIMG]
```

E.g., you have three dimgs in werf.yaml:
```yaml
dimg: app1
from: alpine
---
dimg: app2
from: alpine
---
dimg: app3
from: alpine
```

To build the only `app2` you should use `werf dimg build app2`.


[Q: How can I find image name after build?](#building-4){:id="building-4"}

Use `werf dimg stage image` command for getting image name of last stage or specific stage (`--stage <stage_name>`):

```bash
$ werf dimg stage image
dimgstage-werf:2e29ea2634a335d4e72c801edb55d610cb8657fdf8f77e7257c4b059d2d36e5a

$ werf dimg stage image --stage install
dimgstage-werf:f18fa53176ad78e4dc169e2428c03d79d1e9dde90de1a1890140c3cfdcc33025
```

Or tag your image `werf dimg tag`:

```bash
$ werf dimg tag hello-world
testing_werf
  testing_werf: calculating stages signatures         [RUNNING]
  testing_werf: calculating stages signatures              [OK] 0.5 sec
  custom
    hello-world/testing_werf:latest                   [EXPORTING]
    hello-world/testing_werf:latest                          [OK] 2.11 sec
Running time 2.64 seconds

$ werf dimg tag hello-world --tag-plain test1
testing_werf
  testing_werf: calculating stages signatures         [RUNNING]
  testing_werf: calculating stages signatures              [OK] 0.39 sec
  custom
    hello-world/testing_werf:test1                    [EXPORTING]
    hello-world/testing_werf:test1                           [OK] 2.34 sec
Running time 2.77 seconds
```

[Q: I use werf 0.7, alpine image and get error **standard_init_linux.go:178: exec user process caused "no such file or directory"** on build](#building-5){:id="building-5"}

Werf 0.7 doesn't support `alpine` image, so please use latest werf version.

[Q: Why werf doesn't save cache on failed builds by default?](#building-6){:id="building-6"}

Saving cache on failed builds may cause an incorrect cache. Explanation [here...]({{ site.baseurl }}/not_used/cache_for_advanced_build.html#почему-werf-не-сохраняет-кэш-ошибочных-сборок-по-умолчанию)

## Pushing

[Q: Can I push image to private repository?](#pushing-1){:id="pushing-1"}

Yes, you can use `--registry-username` and `--registry-password` options.

In general for authorization in registry werf use:
* `--registry-username` and `--registry-password` options if specified.
* `CI_JOB_TOKEN` (in CI environment, e.g. GitLab).
* Docker config of a user running werf, `~/.docker/config.json`.


[Q: How can I push to GCR?](#pushing-2){:id="pushing-2"}

To push to GCR you can use the following workaround:

{% raw %}
```bash
werf dimg build
werf dimg tag --tag-ci <REPO>
docker login <REPO>
docker push $(docker images <REPO> --format "{{.Repository}}:{{.Tag}}")
werf dimg flush local
```
{% endraw %}

Werf support push to public and private repositories, but it doesn't work with some platforms such as GCR.


[Q: Can I use several tags at the same time?](#pushing-3){:id="pushing-3"}

Yes.

```bash
werf dimg push --tag custom1 --tag custom2 --tag-build-id --tag-ci --tag-branch --tag-commit
```


## Deploying

[Q: How to debug **Error: render error in...**](#deploying-1){:id="deploying-1"}

You can use `werf kube render` to render templates and `werf kube lint` to validate that it follows the conventions and requirements of the Helm chart standard.


[Q: How to resolve **ErrImagePull** after deploy?](#deploying-2){:id="deploying-2"}

It's not a werf problem. Most likely there is no access to your private repository and if so, you can read about how to add a registry secret in kubernetes documentation [here...](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/).


[Q: How to change helm release name?](#deploying-3){:id="deploying-3"}

Use WERF_HELM_RELEASE_NAME environment.


[Q: How to deploy several applications with different names?](#deploying-4){:id="deploying-4"}

You can pass a variable, e.g. `werf kube deploy --set global.ciProjectName=$CI_PROJECT_NAME ...` and use it in deployment template as {% raw %}`{{ .Values.global.ciProjectName }}`{% endraw %}.