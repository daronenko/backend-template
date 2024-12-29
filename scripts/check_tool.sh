#!/bin/bash

check_and_install() {
    local binary="$1"
    local package="$2"

    if ! command -v "$binary" &> /dev/null; then
        echo "$binary is not installed."
        read -p "Do you want to install it? [Y/n]: " choice
        choice=${choice:-Y}

        if [[ "$choice" =~ ^[Yy]$ ]]; then
            echo "Installing $binary from $package..."
            if go install "$package"; then
                echo "$binary successfully installed."
            else
                echo "Failed to install $binary. Please check your Go environment or permissions."
                exit 1
            fi
        else
            echo "You chose not to install $binary. Exiting..."
            exit 1
        fi
    fi
}

check_and_install "$1" "$2"
