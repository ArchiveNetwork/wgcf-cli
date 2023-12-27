#!/bin/bash
set -e
if [[ ! -z $PREFIX ]]; then
    [[ ! -d $PREFIX ]] && echo -e "error: $PREFIX does not exist." && exit 1
    PREFIX=$(realpath "$PREFIX")
    if [[ -d "$PREFIX/tmp" ]]; then
        WAS_EXIST=true
    else
        mkdir -p $PREFIX/tmp
    fi
    cd $PREFIX/tmp
else
    WAS_EXIST=true
    cd /tmp
fi
OS=$(uname)
if [[ "$OS" != 'Linux' ]] && [[ "$OS" != 'Darwin' ]]; then
    echo "error: This operating system is not supported."
    exit 1
fi

if [[ "$OS" == 'Darwin' ]]; then
    OS='macos'
elif [[ "$OS" == 'Linux' ]]; then
    OS='linux'
fi
if [ -e "/system/bin/app_process" ]; then
        OS='android'
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
        grep Features /proc/cpuinfo | grep -qw 'vfp' || (echo "error: The architecture is not supported." && exit 1)
        ;;
    'armv8' | 'aarch64')
        MACHINE='arm64-v8a'
        ;;
    *)
        echo "error: The architecture is not supported."
        exit 1
        ;;
esac
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
elif [[ "$(type -P brew)" ]]; then
    PACKAGE_MANAGEMENT_INSTALL='brew install'
    PACKAGE_MANAGEMENT_REMOVE='brew uninstall'
    package_provide_tput='ncurses'
fi
if [[ -z $PREFIX ]] && [[ $EUID -ne 0 ]]; then
    echo -e "error: You have to use root to run this script"
    exit 1
elif ! ([ -x $PREFIX ] && [ -w $PREFIX ]);then 
    echo -e "error: You don't have permission to write $PREFIX"
    exit 1
fi
install_software() {
    package_name="$1"
    file_to_detect="$2"
    type -P "$file_to_detect" > /dev/null 2>&1 && return 
    if [[ -z $PACKAGE_MANAGEMENT_INSTALL ]]; then
        echo -e "error: The script does not support the package manager in this operating system."
        echo -e "error: Please install $package_name manually."
        exit 1
    fi
    echo -e "warn: $package_name not installed, installing." && sleep 1
    if ${PACKAGE_MANAGEMENT_INSTALL} "$package_name"; then
        echo "error: $package_name is installed."
    else
        echo -e "error: Installation of $package_name failed, please check your network."
        exit 1
    fi
}
install_software 'unzip' 'unzip'
ZIP_FILE="wgcf-cli.zip"
curl -Lo $ZIP_FILE https://github.com/ArchiveNetwork/wgcf-cli/releases/latest/download/wgcf-cli-${OS}-${MACHINE}.zip
curl -Lo $ZIP_FILE.dgst https://github.com/ArchiveNetwork/wgcf-cli/releases/latest/download/wgcf-cli-${OS}-${MACHINE}.zip.dgst
if [[ $(<$ZIP_FILE.dgst) == 'Not Found' ]]; then
    echo -e "error: No such a version or machine type in release."
    exit 1
fi
CHECKSUM=$(awk -F '= ' '/256=/ {print $2}' "${ZIP_FILE}.dgst")
[[ $OS == 'macos' ]] && LOCALSUM=$((shasum -a 256 wgcf-cli.zip | awk '{printf $1}')|| (sha256sum wgcf-cli.zip | awk '{printf $1}'))
([[ $OS == 'linux' ]] || [[ $OS == android ]]) && LOCALSUM=$(sha256sum "$ZIP_FILE" | awk '{printf $1}')
if [[ "$CHECKSUM" != "$LOCALSUM" ]]; then
    echo -e "error: Checksums do not match, please check your network."
    exit 1
fi
yes|unzip -q $PREFIX/tmp/$ZIP_FILE -d $PREFIX/tmp/wgcf-cli-$OS-$MACHINE > /dev/null 2>&1
if [[ -z $PREFIX ]]; then
    mv /tmp/wgcf-cli-$OS-$MACHINE/wgcf-cli /usr/local/bin/wgcf-cli
else
    if [[ ! -d $PREFIX/bin ]]; then
        mv $PREFIX/tmp/wgcf-cli-$OS-$MACHINE/wgcf-cli $PREFIX/wgcf-cli
    else
        mv $PREFIX/tmp/wgcf-cli-$OS-$MACHINE/wgcf-cli $PREFIX/bin/wgcf-cli
    fi
fi
rm -rf $PREFIX/tmp/wgcf-cli-$OS-$MACHINE/
rm -f $ZIP_FILE
rm -f $ZIP_FILE.dgst
[[ -z $WAS_EXIST ]] && rm -rf $PREFIX/tmp
echo -e "info: wgcf-cli installed."