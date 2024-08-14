#!/bin/bash

# Function to execute the build_and_push script in a service directory
execute_script() {
    local folder=$1
    local script_name="build_and_push.sh"

    echo "Executing build script in $folder"
    cd $folder || { echo "Failed to switch to directory $folder"; exit 1; }

    if [ -f $script_name ]; then
        bash $script_name
        if [ $? -ne 0 ]; then
            echo "Failed to execute script $script_name in $folder."
            exit 1
        fi
    else
        echo "Script $script_name not found in $folder."
        exit 1
    fi

    echo "Switching back to root directory"
    cd - || { echo "Failed to switch back to root directory"; exit 1; }
}

# Execute the build scripts for each service
execute_script "./api-service"
execute_script "./card-quizzler-service"
execute_script "./user-service"
execute_script "./mail-service"

echo "All images have been successfully built and pushed."
