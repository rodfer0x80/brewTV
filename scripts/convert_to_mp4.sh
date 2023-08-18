#!/bin/sh

LIBRARY="/opt/brewTV/library"

for input_file in "$LIBRARY"/*.mkv; do
    if [ -f "$input_file" ]; then
        output_file="${LIBRARY}/$(basename "$input_file" .mkv).mp4"
        ffmpeg -i "$input_file" -c:v libx264 -c:a aac -strict experimental "$output_file"
    fi
done

for input_file in "$LIBRARY"/*.webm; do
    if [ -f "$input_file" ]; then
        output_file="${LIBRARY}/$(basename "$input_file" .mkv).mp4"
        ffmpeg -i "$input_file" -c:v libx264 -c:a aac -strict experimental "$output_file"
    fi
done

for input_file in "$LIBRARY"/*.mp3; do
    if [ -f "$input_file" ]; then
        output_file="${LIBRARY}/$(basename "$input_file" .mkv).mp4"
        ffmpeg -i "$input_file" -c:v libx264 -c:a aac -strict experimental "$output_file"
    fi
done

rm -rf *.mkv *.mp3 *.webm
