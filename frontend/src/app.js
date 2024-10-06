// frontend/src/App.js
import React from 'react';
import ReactPlayer from 'react-player';
import './app.css';

function App() {
  // Construct the base URL based on the current location
  const baseURL = `${window.location.protocol}//${window.location.host}`;
  
  // Video URL
  const videoUrl = `${baseURL}/video`;
  
  // Subtitle URL
  const subtitleUrl = `${baseURL}/sub`;

  return (
    <div className="App">
      <h1>Video Streaming App</h1>
      <ReactPlayer
        url={videoUrl}
        controls
        width="800px"
        height="450px"
        config={{
          file: {
            tracks: [
              {
                kind: 'subtitles',
                src: subtitleUrl,
                srcLang: 'fa',
                default: true,
              },
            ],
          },
        }}
      />
    </div>
  );
}

export default App;
