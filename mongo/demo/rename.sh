#!/bin/bash

basePath=$1
oldName=$2
newName=$3

function exit_1(){
    echo ""
    echo "usage: ./rename.sh <basePath> <oldName> <newName>"
    echo "   eg: ./rename.sh ./ my-app-name wq-account"
    exit 1
}

if [ ! -n "$basePath" ]; then
    echo "error: missing parameter 'basePath'"
    exit_1
fi

if [ ! -n "$oldName" ]; then
    echo "error: missing parameter 'oldName'"
    exit_1
fi

if [ ! -n "$newName" ]; then
    echo "error: missing parameter 'newName'"
    exit_1
fi

# 把第一个字母转为大写
oldName=${oldName^}
newName=${newName^}

function listFiles(){
    cd $1
    items=$(ls)

    for item in $items
    do  
        if [ -d "$item" ]; then
            #echo "this is folder './$item'"
            listFiles $item
        else
            if [ "$item" != "rename.sh" ];then
                # 把第一个字母转为小写
                oldName2=${oldName,}
                newName2=${newName,}

                # 修改文件内容(除了完全匹配字符，也匹配第一个字符为小写字符串)
                sed -i "s/$oldName/$newName/g" $item
                sed -i "s/$oldName2/$newName2/g" $item

                # 修改文件名(规定文件名的第一个字母为小写)
                oldFileName="$item"
                newFileName=$(echo ${oldFileName/"$oldName2"/"$newName2"})
                if [ "$oldFileName" != "$newFileName" ]; then
                    mv $oldFileName $newFileName
                fi
            fi
        fi  
    done 
    cd ..
}

listFiles $basePath

echo "Already replaced the string [$2] to [$3] in the file."

