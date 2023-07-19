#!/bin/sh

CWD="$(dirname $0)/movies"

for input_file in "$CWD"/*.mkv; do
    if [ -f "$input_file" ]; then
        output_file="${CWD}/$(basename "$input_file" .mkv).mp4"
        ffmpeg -i "$input_file" -c:v libx264 -c:a aac -strict experimental "$output_file"
    fi
done

for input_file in "$CWD"/*.webm; do
    if [ -f "$input_file" ]; then
        output_file="${CWD}/$(basename "$input_file" .mkv).mp4"
        ffmpeg -i "$input_file" -c:v libx264 -c:a aac -strict experimental "$output_file"
    fi
done

for input_file in "$CWD"/*.mp3; do
    if [ -f "$input_file" ]; then
        output_file="${CWD}/$(basename "$input_file" .mkv).mp4"
        ffmpeg -i "$input_file" -c:v libx264 -c:a aac -strict experimental "$output_file"
    fi
done
