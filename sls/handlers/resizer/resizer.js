import sharp from 'sharp'
import aws from 'aws-sdk'
const s3 = new aws.S3()

const maxSize = 1200
const jpegQuality = 90

export const handler = async (event) => {
  try {
    const request = event.Records[0].cf.request
    const bucket = request.origin.s3.domainName.split(".")[0]

    const key = decodeURIComponent(request.uri).substring(1)
    const params = {
      Bucket: bucket,
      Key: key
    }

    const data = await s3.getObject(params).promise()
    const sharpBody = sharp(data.Body)
    const metadata = sharpBody.metadata()

    const buffer = await sharpBody
        .jpeg({quality: jpegQuality})
        .resize(maxSize, maxSize, {fit: 'inside'})
        .rotate() // auto-rotated using EXIF Orientation tag
        .toBuffer()

    return {
      status: 200,
      headers: {'content-type': [{key: 'Content-Type', value: `image/${metadata.format}`}]},
      body: buffer.toString('base64'),
      bodyEncoding: 'base64'
    }
  } catch (err) {
    return {
      status: 400,
      body: err.toString(),
    }
  }
}
