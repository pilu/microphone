package main

import (
  "net/http"
  "encoding/json"
)

func RecordingHandler(w http.ResponseWriter, r *http.Request) {
  gid, err := ExtractGidFromRequest(r)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    responseError := ResponseError{ err.Error() }
    json.NewEncoder(w).Encode(responseError)

    return
  }

  searchApiKey, displayApiKey, err := ExtractApAiKeysFromRequest(r)
  if err != nil {
    w.WriteHeader(http.StatusUnauthorized)
    responseError := ResponseError{ err.Error() }
    json.NewEncoder(w).Encode(responseError)
    return
  }

  lyricfindClient.SearchApiKey = searchApiKey
  lyricfindClient.DisplayApiKey = displayApiKey

  recording, err := FindRecording(DB, gid)

  if err != nil && recording != nil {
    w.WriteHeader(http.StatusInternalServerError)
    responseError := ResponseError{ err.Error() }
    json.NewEncoder(w).Encode(responseError)
    return
  }

  if recording == nil {
    err := ResponseError{ "recording not found" }
    w.WriteHeader(http.StatusNotFound)
    responseError := ResponseError{ err.Message }
    json.NewEncoder(w).Encode(responseError)
    return
  }

  userAgent := r.UserAgent()

  lyricsResponse, err := lyricfindClient.SearchAndGetLyrics(recording.Artist, recording.Track, userAgent)

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    responseError := ResponseError{ err.Error() }
    json.NewEncoder(w).Encode(responseError)
    return
  }

  response := BuildRecordingResponse(lyricsResponse)
  json.NewEncoder(w).Encode(response)
}

