<template>
  <h1>Hello</h1>
  <video ref="localVideo" playsinline="" autoplay="" muted=""></video>
  <video ref="remoteVideo" playsinline="" autoplay=""></video>
</template>

<script>

export default {
  name: 'App',
  data() {
    return {
      websocket: null,
      localStream: null
    }
  },
  mounted() {
    console.log(navigator)
    console.log(navigator.mediaDevices)

    navigator.mediaDevices.getUserMedia({audio: true, video: true}).then((stream) => {
      console.log("getUserMedia: ", stream)
    })  .catch((err) => {
      console.error("getUserMedia err: ", err)
    });


    // const room = location.href.substring(location.href.lastIndexOf('/') + 1)
    // console.log("room:", room);
    // const url = "ws://localhost:8080/" + room
    // console.log("connect:", url);
    // this.websocket = new WebSocket(url)
    // this.websocket.addEventListener('open', () => {
    //   // this.websocket.send('Hello Server!');
    // });
    //
    // this.websocket.addEventListener('message', (event) => {
    //   // console.log('Message from server ', event.data);
    // })
  },

  createPeerConnection() {
    const pc = new RTCPeerConnection();
    pc.onicecandidate = e => {
      const message = {
        type: 'candidate',
        candidate: null,
      };
      if (e.candidate) {
        message.candidate = e.candidate.candidate;
        message.sdpMid = e.candidate.sdpMid;
        message.sdpMLineIndex = e.candidate.sdpMLineIndex;
      }
      // signaling.postMessage(message);
      this.websocket.send(message);
    };
    // pc.ontrack = e => remoteVideo.srcObject = e.streams[0];
    // localStream.getTracks().forEach(track => pc.addTrack(track, localStream));
  }
}
</script>

<style>
</style>
