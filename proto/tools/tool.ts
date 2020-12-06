import * as protobuf from "protobufjs";

let filename = process.argv[2]
protobuf.load(`proto/${filename}.proto`).then((root) => {
    let data = JSON.stringify(root.toJSON(), null, 2)
    console.log(data)
});
