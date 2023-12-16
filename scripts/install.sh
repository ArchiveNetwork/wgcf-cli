#!/bin/bash
set -e
cd /tmp
if [[ "$(uname)" != 'Linux' ]]; then
    echo "error: This operating system is not supported."
    exit 1
fi
case "$(uname -m)" in
    'i386' | 'i686')
        MACHINE='32'
    ;;
    'amd64' | 'x86_64')
        MACHINE='64'
    ;;
    'armv7' | 'armv7l')
        MACHINE='arm32-v7a'
        grep Features /proc/cpuinfo | grep -qw 'vfp' || MACHINE='arm32-v5'
    ;;
    'armv8' | 'aarch64')
        MACHINE='arm64-v8a'
    ;;
  *)
    echo "error: The architecture is not supported."
    exit 1
    ;;
esac

package_manager() {
    if [[ "$(type -P apt)" ]]; then
        PACKAGE_MANAGEMENT_INSTALL='apt -y --no-install-recommends install'
        PACKAGE_MANAGEMENT_REMOVE='apt purge'
        package_provide_tput='ncurses-bin'
    elif [[ "$(type -P dnf)" ]]; then
        PACKAGE_MANAGEMENT_INSTALL='dnf -y install'
        PACKAGE_MANAGEMENT_REMOVE='dnf remove'
        package_provide_tput='ncurses'
    elif [[ "$(type -P yum)" ]]; then
        PACKAGE_MANAGEMENT_INSTALL='yum -y install'
        PACKAGE_MANAGEMENT_REMOVE='yum remove'
        package_provide_tput='ncurses'
    elif [[ "$(type -P zypper)" ]]; then
        PACKAGE_MANAGEMENT_INSTALL='zypper install -y --no-recommends'
        PACKAGE_MANAGEMENT_REMOVE='zypper remove'
        package_provide_tput='ncurses-utils'
    elif [[ "$(type -P pacman)" ]]; then
        PACKAGE_MANAGEMENT_INSTALL='pacman -Syu --noconfirm'
        PACKAGE_MANAGEMENT_REMOVE='pacman -Rsn'
        package_provide_tput='ncurses'
     elif [[ "$(type -P emerge)" ]]; then
        PACKAGE_MANAGEMENT_INSTALL='emerge -qv'
        PACKAGE_MANAGEMENT_REMOVE='emerge -Cv'
        package_provide_tput='ncurses'
    else
        echo -e "${ERROR}ERROR:${END} The script does not support the package manager in this operating system."
        exit 1
    fi
}
if [[ $EUID -ne 0 ]]; then
    echo -e "ERROR: You have to use root to run this script"
    exit 1
fi
install_software() {
    package_name="$1"
    file_to_detect="$2"
    type -P "$file_to_detect" > /dev/null 2>&1 && return || echo -e "WARN: $package_name not installed, installing." && sleep 1
    if ${PACKAGE_MANAGEMENT_INSTALL} "$package_name"; then
        echo "INFO: $package_name is installed."
    else
        echo -e "ERROR: Installation of $package_name failed, please check your network."
        exit 1
    fi
}
package_manager
install_software 'unzip' 'unzip'
ZIP_FILE="wgcf-cli.zip"
curl -Lo $ZIP_FILE https://github.com/chise0713/wgcf-cli/releases/latest/download/wgcf-cli-linux-${MACHINE}.zip
curl -Lo $ZIP_FILE.dgst https://github.com/chise0713/wgcf-cli/releases/latest/download/wgcf-cli-linux-${MACHINE}.zip.dgst
CHECKSUM=$(awk -F '= ' '/256=/ {print $2}' "${ZIP_FILE}.dgst")
LOCALSUM=$(sha256sum "$ZIP_FILE" | awk '{printf $1}')
unzip -q $ZIP_FILE -d /tmp/wgcf-cli
mv /tmp/wgcf-cli/wgcf-cli /usr/local/bin/wgcf-cli