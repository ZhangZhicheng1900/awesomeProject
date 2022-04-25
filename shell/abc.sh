#!/bin/bash

env_c="kubelet docker"

declare -a arr_abc=(${env_c})
echo ${arr_abc[0]}