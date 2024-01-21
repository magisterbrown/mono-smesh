/**
 * @file libavformat muxing API usage example
 * @example mux.c
 *
 * Generate a synthetic audio and video signal and mux them to a media file in
 * any supported libavformat format. The default codecs are used.
 */

#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <math.h>

#include <libavcodec/avcodec.h>
#include <libavformat/avformat.h>
#include <libswscale/swscale.h>
#include <libswresample/swresample.h>

#define STREAM_DURATION   10.0
#define STREAM_FRAME_RATE 25 /* 25 images/s */
#define STREAM_PIX_FMT    AV_PIX_FMT_YUV420P /* default pix_fmt */

// a wrapper around a single output AVStream
typedef struct OutputStream {
    AVStream *st;
    AVCodecContext *enc;

    /* pts of the next frame that will be generated */
    int64_t next_pts;
    int samples_count;

    AVFrame *frame;
    AVFrame *tmp_frame;

    AVPacket *tmp_pkt;

    float t, tincr, tincr2;

    struct SwsContext *sws_ctx;
    struct SwrContext *swr_ctx;
} OutputStream;



static int write_frame(AVFormatContext *fmt_ctx, AVCodecContext *c,
                       AVStream *st, AVFrame *frame, AVPacket *pkt)
{
    int ret;

    // send the frame to the encoder
    ret = avcodec_send_frame(c, frame);
    while (ret >= 0) {
        ret = avcodec_receive_packet(c, pkt);
        if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF)
            break;
        /* rescale output packet timestamp values from codec to stream timebase */
        av_packet_rescale_ts(pkt, c->time_base, st->time_base);
        pkt->stream_index = st->index;

        /* Write the compressed frame to the media file. */
        ret = av_interleaved_write_frame(fmt_ctx, pkt);
        /* pkt is now blank (av_interleaved_write_frame() takes ownership of
         * its contents and resets pkt), so that no unreferencing is necessary.
         * This would be different if one used av_write_frame(). */
       
    }

    return ret == AVERROR_EOF ? 1 : 0;
}

/* Add an output stream. */
static void add_stream(OutputStream *ost, AVFormatContext *oc,
                       const AVCodec **codec,
                       enum AVCodecID codec_id)
{
    AVCodecContext *c;
    /* find the encoder */
    *codec = avcodec_find_encoder(codec_id);
    ost->tmp_pkt = av_packet_alloc();
    ost->st = avformat_new_stream(oc, NULL);
    ost->st->id = oc->nb_streams-1;
    c = avcodec_alloc_context3(*codec);
    ost->enc = c;

    c->codec_id = codec_id;

    c->bit_rate = 400000;
    /* Resolution must be a multiple of two. */
    c->width    = 352;
    c->height   = 288;
    /* timebase: This is the fundamental unit of time (in seconds) in terms
     * of which frame timestamps are represented. For fixed-fps content,
     * timebase should be 1/framerate and timestamp increments should be
     * identical to 1. */
    ost->st->time_base = (AVRational){ 1, STREAM_FRAME_RATE };
    c->time_base       = ost->st->time_base;

    c->gop_size      = 12; /* emit one intra frame every twelve frames at most */
    c->pix_fmt       = STREAM_PIX_FMT;

    /* Some formats want stream headers to be separate. */
    c->flags |= AV_CODEC_FLAG_GLOBAL_HEADER;
}

/**************************************************************/
/* video output */

static void open_video(AVFormatContext *oc, const AVCodec *codec,
                       OutputStream *ost)
{
    AVCodecContext *c = ost->enc;


    /* open the codec */
    avcodec_open2(c, codec, NULL);

    ost->frame = av_frame_alloc();

    ost->frame->format = c->pix_fmt;
    ost->frame->width  = c->width;
    ost->frame->height = c->height;

    /* allocate the buffers for the frame data */
    av_frame_get_buffer(ost->frame, 0);

    /* copy the stream parameters to the muxer */
    avcodec_parameters_from_context(ost->st->codecpar, c);
    
}

/* Prepare a dummy image. */
static void fill_yuv_image(AVFrame *pict, int frame_index,
                           int width, int height)
{
    int x, y, i;

    i = frame_index;

    /* Y */
    for (y = 0; y < height; y++)
        for (x = 0; x < width; x++)
            pict->data[0][y * pict->linesize[0] + x] = x + y + i * 3;

    /* Cb and Cr */
    for (y = 0; y < height / 2; y++) {
        for (x = 0; x < width / 2; x++) {
            pict->data[1][y * pict->linesize[1] + x] = 128 + y + i * 2;
            pict->data[2][y * pict->linesize[2] + x] = 64 + x + i * 5;
        }
    }
}

static AVFrame *get_video_frame(OutputStream *ost)
{
    AVCodecContext *c = ost->enc;

    /* check if we want to generate more frames */
    if (av_compare_ts(ost->next_pts, c->time_base,
                      STREAM_DURATION, (AVRational){ 1, 1 }) > 0)
        return NULL;
   
    fill_yuv_image(ost->frame, ost->next_pts, c->width, c->height);

    ost->frame->pts = ost->next_pts++;

    return ost->frame;
}

static void close_stream(AVFormatContext *oc, OutputStream *ost)
{
    avcodec_free_context(&ost->enc);
    av_frame_free(&ost->frame);
    av_frame_free(&ost->tmp_frame);
    av_packet_free(&ost->tmp_pkt);
    sws_freeContext(ost->sws_ctx);
    swr_free(&ost->swr_ctx);
}

/**************************************************************/
/* media file output */

int main(int argc, char **argv)
{
    OutputStream video_st = { 0 };
    const AVOutputFormat *fmt;
    AVFormatContext *oc;
    const AVCodec *video_codec;
    
    /* allocate the output media context */
    avformat_alloc_output_context2(&oc, NULL, "mp4", NULL);
    fmt = oc->oformat;

    /* Add the audio and video streams using the default format codecs
     * and initialize the codecs. */
    add_stream(&video_st, oc, &video_codec, fmt->video_codec);

    /* Now that all the parameters are set, we can open the audio and
     * video codecs and allocate the necessary encode buffers. */
    open_video(oc, video_codec, &video_st);

    //av_dump_format(oc, 0, filename, 1);

    /* open the output file, if needed */
    //avio_open(&oc->pb, filename, AVIO_FLAG_WRITE);
    avio_open_dyn_buf(&oc->pb);
       /* Write the stream header, if any. */
    avformat_write_header(oc, NULL);
    int encode_video=1;
    while (encode_video) {
        /* select the stream to encode */
        encode_video = !write_frame(oc, video_st.enc, video_st.st, get_video_frame(&video_st), video_st.tmp_pkt); 
    }

    av_write_trailer(oc);

    /* Close each codec. */
    close_stream(oc, &video_st);

    uint8_t *sps;
    int sz = avio_close_dyn_buf(oc->pb, &sps);
    for(uint8_t *i=sps;i<sps+10;i++) {
        printf("Header: %d\n", *i);
    }
    FILE *f;
    f = fopen("videold.mp4", "wb");
    fwrite(sps, 1, sz, f);
    fclose(f);


    return 0;
}
