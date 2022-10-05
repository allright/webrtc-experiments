<template>
  <h1>Hello</h1>
  <video ref="localVideo" playsinline="" autoplay="" muted=""></video>
  <hr/>
  <video ref="remoteVideo" playsinline="" autoplay=""></video>
</template>

<script>

export default {
  name: 'App',
  data() {
    return {
      websocket: null,
      localStream: null,
      pc: null,
      dc: null,
      pingCount: 0,
      nextPingWaitingCount: 0,
      prevRecvTS: new Date().getTime()
    }
  },
  methods: {

    send(msg) {
      console.log("send:", msg)
      this.websocket.send(JSON.stringify(msg))
    },

    async start() {
      this.localStream = await navigator.mediaDevices.getUserMedia({audio: true, video: true});
      this.$refs.localVideo.srcObject = this.localStream;
      console.log("start")
      this.send({"type": "ready"})
    },

    onMessage(event) {
      const msg = JSON.parse(event.data)
      console.log("onMessage:", msg)
      switch (msg.type) {
        case "ready":
          if (this.pc) {
            console.log('already in call, ignoring');
            return;
          }
          this.makeCall()
          break
        case 'offer':
          this.handleOffer(msg)
          break
        case 'answer':
          this.handleAnswer(msg)
          break
        case 'candidate':
          this.handleCandidate(msg)
          break
      }
    },


    // https://developer.mozilla.org/en-US/docs/Web/API/WebRTC_API/Simple_RTCDataChannel_sample
    // DataChannel example
    createPeerConnection() {
      this.pc = new RTCPeerConnection();
      this.pc.onicecandidate = e => {
        const message = {
          type: 'candidate',
          candidate: null,
        };
        if (e.candidate) {
          message.candidate = e.candidate.candidate;
          message.sdpMid = e.candidate.sdpMid;
          message.sdpMLineIndex = e.candidate.sdpMLineIndex;
        }
        this.send(message);
      };
      const remoteVideo = this.$refs.remoteVideo
      this.pc.ontrack = e => remoteVideo.srcObject = e.streams[0];
      this.localStream.getTracks().forEach(track => this.pc.addTrack(track, this.localStream));

      this.pc.ondatachannel = (event) => {
        console.log("on data channel:", event)
        const rcvChannel = event.channel
        rcvChannel.onmessage = (event) => {
       //   console.log("onmessage:", event)
          const time = new Date()
          const msg = JSON.parse(event.data)
          switch (msg.msg) {
            case "ping": {
              const pong = {
                "msg": "pong",
                "ping_id": msg.id,
                "ping_ts": msg.timestamp,
                "timestamp": time.getTime()
              }
          //    console.log("send pong:", pong)
              this.dc.send(JSON.stringify(pong))
            }
              break

            case "pong": {
              const delta_ms = time - msg.ping_ts
              const delta_ping_num = msg.ping_id - this.nextPingWaitingCount
              this.nextPingWaitingCount = msg.ping_id + 1
              const diffTS = msg.ping_ts - this.prevRecvTS
              this.prevRecvTS = msg.ping_ts
              console.log("latency (id:" +  msg.ping_id + ", lost: " + delta_ping_num + ", delta_ping: " + delta_ms + ", diff: " + diffTS + ")")
            }
              break

          }

        }
        rcvChannel.onopen = (event) => {
          console.log("onopen:", event)
        }
        rcvChannel.onclose = (event) => {
          console.log("onopen:", event)
        }
      }

      this.dc = this.pc.createDataChannel("sendChannel")
      this.dc.onopen = () => {
        setInterval(this.sendPing, 10)
      }
    },

    sendPing() {
      const time = new Date()
   //   console.log("ping", time, time.getTime())
      const ping = {
        "msg": "ping",
        "id": this.pingCount++,
        "timestamp": time.getTime()
      }
      this.dc.send(JSON.stringify(ping))
    },

    async makeCall() {
      await this.createPeerConnection();
      const offer = await this.pc.createOffer();
      this.send({type: 'offer', sdp: offer.sdp});
      await this.pc.setLocalDescription(offer);
    },

    async handleOffer(offer) {
      console.log("handleOffer:", offer)
      if (this.pc) {
        console.error('existing peerconnection');
        return;
      }
      await this.createPeerConnection();
      await this.pc.setRemoteDescription(offer);

      const answer = await this.pc.createAnswer();
      this.send({type: 'answer', sdp: answer.sdp});
      await this.pc.setLocalDescription(answer);
    },

    async handleAnswer(answer) {
      if (!this.pc) {
        console.error('no peerconnection');
        return;
      }
      await this.pc.setRemoteDescription(answer);
    },

    async handleCandidate(candidate) {
      if (!this.pc) {
        console.error('no peerconnection');
        return;
      }
      if (!candidate.candidate) {
        await this.pc.addIceCandidate(null);
      } else {
        await this.pc.addIceCandidate(candidate);
      }
    }

  },

  mounted() {
    const room = location.href.substring(location.href.lastIndexOf('/') + 1)
    console.log("room:", room);
    //const url = "wss://localhost:8443/ws/" + room
    console.log(window.location)
    const url = "wss://" + window.location.host + "/ws/" + room
    console.log("connect:", url);
    this.websocket = new WebSocket(url)
    this.websocket.addEventListener('open', () => this.start());
    this.websocket.addEventListener('message', (event) => this.onMessage(event))
  }
}
</script>

<style>
</style>
