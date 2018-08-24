

if kextstat | grep -q "com.qtt.xuzhiqiang.gotproxy"; then
    echo "already loaded"
else
    sudo rm -rf /tmp/darwin_kext.kext && sudo cp -R darwin/darwin_kext.kext /tmp/darwin_kext.kext
    sudo kextload  /tmp/darwin_kext.kext
fi

