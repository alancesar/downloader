package testdata

const RedGIFsGifResponse = `
{
  "gif": {
    "id": "some-media",
    "createDate": 1680462051,
    "hasAudio": true,
    "width": 1920,
    "height": 1080,
    "hls": true,
    "likes": 32,
    "niches": ["some-niche"],
    "tags": ["Some Tag"],
    "verified": true,
    "views": null,
    "description": "Some description",
    "duration": 30,
    "published": true,
    "urls": {
      "thumbnail": "https://thumbs44.redgifs.com/some-media-mobile.jpg?expires=1680625800&signature=v2:755a40e0b354717fc9cfa250cc01695e89ce2b830d9f30973fd5498acace10a1&for=192.198.0.1&hash=6163438793",
      "sd": "https://thumbs44.redgifs.com/some-media-mobile.mp4?expires=1680625800&signature=v2:a7c3be8860e39bba9e69ebc96f3bd288cd9e251647a333ea18e86e28b52c7b95&for=192.198.0.1&hash=6163438793",
      "vthumbnail": "https://thumbs44.redgifs.com/some-media-mobile.mp4?expires=1680625800&signature=v2:a7c3be8860e39bba9e69ebc96f3bd288cd9e251647a333ea18e86e28b52c7b95&for=192.198.0.1&hash=6163438793",
      "hd": "https://thumbs44.redgifs.com/some-media.mp4?expires=1680625800&signature=v2:d330e63fd8fbf4dda3f7d66d399b9c1fffdc7383c1190bee65f694b9f5affb29&for=192.198.0.1&hash=6163438793",
      "poster": "https://thumbs44.redgifs.com/some-media-poster.jpg?expires=1680625800&signature=v2:761f3ab07837f9066aa605237fd3f82fa9f0097b65490ca5af892f40a84b1016&for=192.198.0.1&hash=6163438793"
    },
    "userName": "some-username",
    "type": 1,
    "avgColor": "#000000",
    "gallery": "some-gallery",
    "hideHome": false,
    "hideTrending": false,
    "sexuality": ["straight"]
  },
  "user": null,
  "niches": []
}
`

const RedGIFsUserSearchResponse = `
{
  "page": 1,
  "pages": 1,
  "total": 2,
  "gifs": [
    {
      "avgColor": "#000000",
      "createDate": 1664178721,
      "duration": 30,
      "gallery": "some-gallery-1",
      "hasAudio": true,
      "height": 1080,
      "hideHome": false,
      "hideTrending": false,
      "hls": false,
      "id": "some-media-1",
      "likes": 256,
      "niches": ["some-niche"],
      "published": true,
      "tags": [
        "Some Tag"
      ],
      "type": 1,
      "urls": {
        "sd": "https://thumbs44.redgifs.com/some-media-1-mobile.mp4?expires=1680626400&signature=v2:21f8d929a5df83247f8bbad2c421d1925683f976aef926633f4f96a6c94853fa&for=192.168.0.1&hash=7011125643",
        "hd": "https://thumbs44.redgifs.com/some-media-1.mp4?expires=1680626400&signature=v2:d440f129db5226483ea33df3320bfe61eef74a2598d541ab36de511e1568cb0a&for=192.168.0.1&hash=7011125643",
        "poster": "https://thumbs44.redgifs.com/some-media-1-poster.jpg?expires=1680626400&signature=v2:929b6438a816494bbd520f5a31b480804da4b746683b897c715d8a194acd23d9&for=192.168.0.1&hash=7011125643",
        "thumbnail": "https://thumbs44.redgifs.com/some-media-1-mobile.jpg?expires=1680626400&signature=v2:569d62d0068f1f6f030a060fffb976ff7435a86974e2fd43fa541e81fee2692b&for=192.168.0.1&hash=7011125643",
        "vthumbnail": "https://thumbs44.redgifs.com/some-media-1-mobile.mp4?expires=1680626400&signature=v2:21f8d929a5df83247f8bbad2c421d1925683f976aef926633f4f96a6c94853fa&for=192.168.0.1&hash=7011125643"
      },
      "userName": "some-username",
      "verified": true,
      "views": 1000,
      "width": 1920,
      "sexuality": ["straight"]
    },
    {
      "avgColor": "#000000",
      "createDate": 1670704347,
      "duration": 42,
      "gallery": "some-gallery-2",
      "hasAudio": true,
      "height": 1080,
      "hideHome": false,
      "hideTrending": false,
      "hls": false,
      "id": "some-media-2",
      "likes": 128,
      "niches": ["some-niche"],
      "published": true,
      "tags": [
        "Some Tag"
      ],
      "type": 1,
      "urls": {
        "sd": "https://thumbs44.redgifs.com/some-media-2-mobile.mp4?expires=1680626400&signature=v2:04476e1f2872004d32ef0d67c01b6871b0ac79d09343ba279b52c1d579993264&for=192.168.0.1&hash=7011125643",
        "hd": "https://thumbs44.redgifs.com/some-media-2.mp4?expires=1680626400&signature=v2:12b4a083d11d077dcdf489020a8c4cc8e90f3464fc9927d5eb2cff69efe71001&for=192.168.0.1&hash=7011125643",
        "poster": "https://thumbs44.redgifs.com/some-media-2-poster.jpg?expires=1680626400&signature=v2:0883fd49f94b20c6711f64c9cd0b2b47c5f9c82ad68d296e4c0eb426973d7699&for=192.168.0.1&hash=7011125643",
        "thumbnail": "https://thumbs44.redgifs.com/some-media-2-mobile.jpg?expires=1680626400&signature=v2:13667fe3b8175459613770cc81cbf7db25a75e52c7cfc4d8cb2da8c63a4b79c4&for=192.168.0.1&hash=7011125643",
        "vthumbnail": "https://thumbs44.redgifs.com/some-media-2-mobile.mp4?expires=1680626400&signature=v2:04476e1f2872004d32ef0d67c01b6871b0ac79d09343ba279b52c1d579993264&for=192.168.0.1&hash=7011125643"
      },
      "userName": "some-username",
      "verified": true,
      "views": 1600,
      "width": 720,
      "sexuality": ["straight"]
    }
}
`

const RedGIFsTokenResponse = `
{
    "addr": "192.168.0.1",
    "agent": "downloaddit-v0",
    "rtfm": "https://github.com/Redgifs/api/wiki/Temporary-tokens",
    "token": "some.token.here"
}
`
