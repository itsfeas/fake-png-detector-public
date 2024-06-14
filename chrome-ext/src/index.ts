// import { FancyArray, array, at, base, broadcastArray, broadcastArrays, castingModes, dataBuffer, defaults, dispatch, dtype, dtypes, empty, emptyLike, flag, flags, ind2sub, indexModes, iter, maybeBroadcastArray, maybeBroadcastArrays, minDataType, mostlySafeCasts, ndarray, ndarray2array, ndims, ndsliceAssign, nextDataType, numel, numelDimension, offset, order, orders, outputDataTypePolicies, promotionRules, safeCasts, sameKindCasts, scalar2ndarray, shape, slice, sliceDimension, sliceDimensionFrom, sliceDimensionTo, sliceFrom, sliceTo, stride, strides, sub2ind, zeros, zerosLike } from 'https://cdn.jsdelivr.net/gh/stdlib-js/ndarray@esm/index.mjs';
import * as ort from 'onnxruntime-web';
import ndarray from 'ndarray';
import { subseq, assign, divseq } from 'ndarray-ops';
import * as imgLoader from "./image-loader";

const PROXY_URL = "http://localhost:8081/?url=";
const IMAGE_SIZE = 224;

// Create an ONNX inference session with WebGL backend.
var session: null | ort.InferenceSession = null;
const imageCache = new Map<string, boolean>();

const updateImageOverlay = (overlay: HTMLDivElement, element: HTMLImageElement): (isReal: boolean) => void => {
    return (inference: boolean) => {
        if (!inference) {
            // overlay.style.backgroundColor = 'red';
            overlay.style.opacity = '30%';
            overlay.style.top = element.offsetTop.toString();
            overlay.style.left = element.offsetLeft.toString();
            overlay.style.width = element.width.toString();
            overlay.style.height = element.height.toString();
        }
        console.log("inference", element.src, inference);
    };
}

const createImageOverlay = (element: HTMLImageElement, superOverlay: HTMLDivElement) => {
    const overlay = document.createElement("div");
    overlay.style.position = 'absolute';
    overlay.style.top = element.offsetTop.toString();
    overlay.style.left = element.offsetLeft.toString();
    overlay.style.width = element.width.toString();
    overlay.style.height = element.height.toString();
    overlay.style.backgroundColor = 'red';
    overlay.style.opacity = '0%';
    superOverlay.appendChild(overlay);
    return overlay;
}

function createSuperOverlay() {
    const preExistingOverlay = document.getElementById("fake-png-detector-super-overlay");
    if (preExistingOverlay) {
        // console.log("removed!");
        preExistingOverlay.remove();
    }
    const superOverlay = document.createElement("div");
    document.body.appendChild(superOverlay);
    superOverlay.id = "fake-png-detector-super-overlay";
    superOverlay.style.zIndex = "100";
    superOverlay.style.width = "100%";
    superOverlay.style.height = "100%";
    superOverlay.style.position = 'absolute';
    superOverlay.style.pointerEvents = 'none';
    superOverlay.style.top = "0";
    superOverlay.style.left = "0";
    return superOverlay;
}

// function createOverlayContainer() {
//     const preExistingOverlays = document.getElementsByClassName("fake-png-detector-overlay-container");
//     for (let index = 0; index < preExistingOverlays.length; index++) {
//         preExistingOverlays[index].remove();
//     }
//     const images = document.getElementsByTagName("img");
//     for (let index = 0; index < images.length; index++) {
//         images[index].parentNode
//     }
//     const superOverlay = document.createElement("div");
//     document.body.appendChild(superOverlay);
//     superOverlay.id = "fake-png-detector-super-overlay";
//     superOverlay.style.zIndex = "100";
//     superOverlay.style.width = "100%";
//     superOverlay.style.height = "100%";
//     superOverlay.style.position = 'absolute';
//     superOverlay.style.pointerEvents = 'none';
//     superOverlay.style.top = "0";
//     superOverlay.style.left = "0";
//     return superOverlay;
// }


const inferenceOnImageData = async (imageData: ImageData) => {
    await startSession();
    const preprocessedData = preprocess(imageData.data);
    const inputTensor = new ort.Tensor('float32', preprocessedData, [1, 3, IMAGE_SIZE, IMAGE_SIZE]);
    // console.log("inputTensor", inputTensor);
    const outputMap = await session.run({ modelInput: inputTensor });
    const outputData = outputMap.modelOutput.data;
    // console.log("done!", outputData);
    return outputData as Float32Array;
}


const startSession = async () => {
    if (!session) {
        session = await ort.InferenceSession.create("./SQUEEZE.onnx", { graphOptimizationLevel: 'all', executionMode: 'sequential', intraOpNumThreads: 4 });
    }
}


/**
 * Preprocess raw image data to match SqueezeNet requirement.
 */
const mean = [0.485, 0.456, 0.406];
const std = [0.229, 0.224, 0.225];
const preprocess = (data: Uint8ClampedArray, width: number = IMAGE_SIZE, height: number = IMAGE_SIZE) => {
    const dataFromImage = ndarray(new Float32Array(data), [width, height, 4]);
    const dataProcessed = ndarray(new Float32Array(width * height * 3), [1, 3, height, width]);

    // Normalize 0-255 to (-1)-1
    divseq(dataFromImage, 128.0);
    subseq(dataFromImage, 1.0);
    for (let i = 0; i < 3; i++) {
        subseq(dataFromImage.pick(i, null, null), mean[i]);
        divseq(dataFromImage.pick(i, null, null), std[i]);
    }
    
    // Realign imageData from [224*224*4] to the correct dimension [1*3*224*224].
    assign(dataProcessed.pick(0, 0, null, null), dataFromImage.pick(null, null, 2));
    assign(dataProcessed.pick(0, 1, null, null), dataFromImage.pick(null, null, 1));
    assign(dataProcessed.pick(0, 2, null, null), dataFromImage.pick(null, null, 0));
    // assign(dataProcessed.pick(0, 0, null, null), dataFromImage.pick(null, null, 0));
    // assign(dataProcessed.pick(0, 1, null, null), dataFromImage.pick(null, null, 1));
    // assign(dataProcessed.pick(0, 2, null, null), dataFromImage.pick(null, null, 2));
    return dataProcessed.data;
}

const init = async () => {
    const images = document.getElementsByTagName('img');
    const superOverlay = createSuperOverlay();
    for (let index = 0; index < images.length; index++) {
        const element = images[index];
        const overlay = createImageOverlay(element, superOverlay);
        if (element.complete) {
            isImageReal(element, updateImageOverlay(overlay, element));
        } else {
            element.addEventListener("load", (e) => {
                if (e.target === undefined) {
                    console.log("undefined", element.src)
                    return;
                }
                const img = e.target as HTMLImageElement;
                console.log("late load", img);
                isImageReal(img, updateImageOverlay(overlay, img));
                // const imgData = (element.src.includes("localhost") || element.src.includes("data:image/")) ? await imgLoader.loadImageFromImageElement(element) : await imgLoader.loadImageFromURL(imgUrl);
                // const inference = await inferenceOnImageData(imgData);
                // console.log("inference", element.src, inference);
            }, false);
        }
    }
}

const isImageReal = async (element: HTMLImageElement, callback?: (isReal: boolean) => void) => {
    let ret = false;
    if (imageCache.has(element.src)) {
        ret = imageCache.get(element.src);
    } else {
        const imgUrl = PROXY_URL + encodeURIComponent(element.src);
        const imgData = (element.src.includes("localhost") || element.src.includes("data:image/")) ? await imgLoader.loadImageFromImageElement(element) : await imgLoader.loadImageFromURL(imgUrl);
        const output = await inferenceOnImageData(imgData);
        ret = output[1]>output[0] ? true : false;
        // imageCache.set(element.src, ret);
    }
    if (callback) callback(ret);
    return ret;
}

// startSession();
init();

window.addEventListener("resize", (e) => init(), true);