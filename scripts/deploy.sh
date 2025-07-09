#!/bin/bash
plugin_name="csi.san.synology.com"

SCRIPT_PATH="$(realpath "$0")"
SOURCE_PATH="$(realpath "$(dirname "${SCRIPT_PATH}")"/../)"
config_file="${SOURCE_PATH}/config/client-info.yml"
plugin_dir="/var/lib/kubelet/plugins/$plugin_name"

source "$SOURCE_PATH"/scripts/functions.sh

# 1. Build
csi_build(){
    echo "==== Build synology-csi .... ===="
    source "$SOURCE_PATH"/build.sh
}

# 2. Install
csi_install(){
    echo "==== Creates namespace and secrets, then installs synology-csi ===="
    parse_version
    echo "Deploy Version: $deploy_k8s_version"

    kubectl create ns synology-csi
    kubectl create secret -n synology-csi generic client-info-secret --from-file="$config_file"

    if [ ! -d "$plugin_dir" ]; then
        mkdir -p $plugin_dir
    fi

    kubectl apply -f "$SOURCE_PATH"/deploy/kubernetes/$deploy_k8s_version

    if [ "$basic_mode" == false ]; then
        kubectl apply -f "$SOURCE_PATH"/deploy/kubernetes/$deploy_k8s_version/snapshotter
    fi

    if [ "$openshift_mode" == true ]; then
        kubectl apply -f "$SOURCE_PATH"/deploy/kubernetes/openshift_synology_scc.yml
    fi

    if [ "$talos_mode" == true ]; then
        kubectl patch daemonset \
            synology-csi-node \
            --namespace synology-csi \
            --type='json' \
            -p='[{"op": "replace", "path": "/spec/template/spec/containers/1/args", "value": [
            "--nodeid=$(KUBE_NODE_NAME)",
            "--endpoint=$(CSI_ENDPOINT)",
            "--client-info",
            "/etc/synology/client-info.yml",
            "--log-level=info",
            "--iscsiadm-path=/usr/local/sbin/iscsiadm"
        ]}]'
    fi
}

# 3. Uninstall
csi_uninstall(){
    kubectl delete -f "$SOURCE_PATH"/deploy/kubernetes/$deploy_k8s_version
    kubectl delete -f "$SOURCE_PATH"/deploy/kubernetes/$deploy_k8s_version/snapshotter
    kubectl delete -f "$SOURCE_PATH"/deploy/kubernetes/openshift_synology_scc.yml
}

print_usage(){
    echo "Usage:"
    echo "    deploy.sh run          build and install"
    echo "    deploy.sh [command]    specify an action to be performed"
    echo "Available Commands:"
    echo "    build                  build docker image only"
    echo "    install [flag]         install csi plugin with the specified flag"
    echo "    uninstall              uninstall the csi plugin and snapshot controller"
    echo "    help                   show help"
    echo "Available Flags:"
    echo "    -a, --all              deploy csi plugin and snapshotter"
    echo "    -b, --basic            deploy basic csi plugin only"
    echo "    -o, --openshift        deploy on openshift cluster"
    echo "    -t, --talos            deploy on talos cluster"
    echo "Examples:"
    echo "    deploy.sh run"
    echo "    deploy.sh install --basic"
}

basic_mode=false
parse_flags(){
    case "$1" in
        -a|--all)
            ;;
        -b|--basic)
            basic_mode=true
            ;;
        -o|--openshift)
            openshift_mode=true
            ;;
        -t|--talos)
            talos_mode=true
            ;;
        *)
            print_usage
            exit 1
            ;;
    esac
}

case "$1" in
    build)
        csi_build
        ;;
    install)
        parse_flags "$2"
        csi_install
        ;;
    run)
        csi_build
        csi_install
        ;;
    uninstall)
        csi_uninstall
        ;;
    *)
        print_usage
        exit 1
        ;;
esac
