#!/bin/sh

convert_all_files_in_dir_by_format() {
    # Convert all files in a directory by format
    # Usage: convert_all_files_in_dir_by_format <WORKDIR> <INPUT_FILE_FORMAT> <OUTPUT_FILE_FORMAT>
    # Example: convert_all_files_in_dir_by_format "/opt/brewTV/library" "webm" "mp4"
    local WORKDIR="$1"
    local INPUT_FILE_FORMAT="$2"
    local OUTPUT_FILE_FORMAT="$3"
    
    for input_file in "$WORKDIR"/*."$FILE_FORMAT"; do
        if [ -f "$input_file" ]; then
            output_file="${WORKDIR}/$(basename "$input_file" ".$FILE_FORMAT").$OUTPUT_FILE_FORMAT"
            ffmpeg -i "$input_file" -c:v libx264 -c:a aac -strict experimental "$output_file"
        fi
    done
    rm -rf "$WORKDIR/*.$FILE_FORMAT"
}