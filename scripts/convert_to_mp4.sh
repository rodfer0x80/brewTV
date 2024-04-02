#!/bin/bash

if ! command -v ffmpeg &> /dev/null
then
    echo "ffmpeg is not installed. Please install ffmpeg first."
    exit
fi

print_usage() {
    echo "Usage: $0 [-d] <directory>"
}

if [ $# -eq 0 ]; then
    print_usage
    exit 1
fi

delete_original=false
while getopts ":d" opt; do
    case ${opt} in
        d )
            delete_original=true
            ;;
        \? )
            echo "Invalid option: -$OPTARG" 1>&2
            print_usage
            exit 1
            ;;
        : )
            echo "Option -$OPTARG requires an argument." 1>&2
            print_usage
            exit 1
            ;;
    esac
done
shift $((OPTIND -1))

if [ ! -d "$1" ]; then
    echo "Directory '$1' does not exist."
    exit 1
fi

cd "$1" || exit
mkdir -p original
mv * original/
for file in original/*; do
    extension="${file##*.}"
    if [ "$extension" != "mp4" ]; then
        output_file="${file%.*}.mp4"
        ffmpeg -i "$file" "$output_file"
        mv "$output_file" .
    fi
done

if [ "$delete_original" = true ]; then
    rm -rf original
    echo "Original directory deleted"
fi

echo "Conversion completed"

