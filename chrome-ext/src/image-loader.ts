import loadImage from 'blueimp-load-image';

const IMG_SIZE = 224;

export const loadImageFromURL = async (url: string) => {
    const loaded = await loadImage(url, { canvas: true, crossOrigin: "*" });
    const canvas = loaded.image as HTMLCanvasElement;
    const context = canvas.getContext('2d');
    const imageData = context.getImageData(0, 0, IMG_SIZE, IMG_SIZE);
    return imageData;
}

export const loadImageFromImageElement = async (img: HTMLImageElement) => {
    const loaded = await loadImage(img.src, { canvas: true });
    const canvas = loaded.image as HTMLCanvasElement;
    const context = canvas.getContext('2d');
    context.drawImage(img, 0, 0);
    return context.getImageData(0, 0, IMG_SIZE, IMG_SIZE);
}

export const loadImageFromBlob = async (blob: Blob) => {
    // const url = URL.createObjectURL(b);
    // const img = new Image();

    // img.onload = (e) => {
    //     if (img) {
    //         URL.revokeObjectURL(e.target.src);             // free memory held by Object URL
    //         c.getContext("2d").drawImage(this, 0, 0);  // draw image onto canvas (lazy method™)
    //     }
    // };

    // img.src = url;  

    // const loaded = await loadImage(b, { canvas: true });
    // const canvas = loaded.image as HTMLCanvasElement;
    // const context = canvas.getContext('2d');
    // context.drawImage(img, 0, 0);
    // return context.getImageData(0, 0, img.width, img.height);
    let blobUrl = URL.createObjectURL(blob);

    return new Promise((resolve, reject) => {
        let img = new Image();
        img.onload = () => resolve(img);
        img.onerror = err => reject(err);
        img.src = blobUrl;
    }).then((img: HTMLImageElement) => {
        URL.revokeObjectURL(blobUrl);
        // Limit to 256x256px while preserving aspect ratio
        let [w, h] = [img.width, img.height]
        let aspectRatio = w / h
        // Say the file is 1920x1080
        // divide max(w,h) by 256 to get factor
        let factor = Math.max(w, h) / 256
        w = w / factor
        h = h / factor

        // REMINDER
        // 256x256 = 65536 pixels with 4 channels (RGBA) = 262144 data points for each image
        // Data is encoded as Uint8ClampedArray with BYTES_PER_ELEMENT = 1
        // So each images = 262144bytes
        // 1000 images = 260Mb
        let canvas = document.createElement("canvas");
        canvas.width = w;
        canvas.height = h;
        let ctx = canvas.getContext("2d");
        ctx.drawImage(img, 0, 0);

        return ctx.getImageData(0, 0, w, h);    // some browsers synchronously decode image here
    })
}