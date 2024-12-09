import React, { memo ,useRef,useEffect,useState} from 'react';
import { Editor } from '@tinymce/tinymce-react';
import { myGlobalInfoStore } from '@/stores';

function EditorTinyMCE ({
        editorPlaceholder = '',
        className = '',
        // initialValue='',
       value='',
        onChange,
        // onFocus,
        // onBlur,
        // autoFocus = false,
    })
{
    const [contentValue, setContentValue] = useState(value ?? '');
    useEffect(() => setContentValue(value ?? ''), [value]);

    const {isSideNavSticky,sideNavStickyTop}= myGlobalInfoStore()

    console.log("editorTinyMCE value:")
    const tinyMceKey = 'pm4bf4u8cw7y3w24vo5vrwmh09tgj9qcgns63w0293niwzpk'
    const templateStr="写点什么吧";

    const toolbar=[
        // ' blocks styles fontfamily  fontsize   | pastetext ',
        ' blocks  fontfamily  fontsize   | code pastetext ',
       
        'formatpainter forecolor backcolor bold italic underline strikethrough link anchor | alignleft aligncenter alignright alignjustify outdent indent |  bullist numlist |  removeformat | table image media  emoticons hr    preview | fullscreen | bdmap indent2em lineheight  axupimgs',
        ' undo redo restoredraft ',
        ' blockquote subscript superscript ',
        //不要的 charmap 特殊字符
        //分页符 pagebreak
        //insertdatetime 时间日期
        //print 打印


    ];
    const toolBarStr=toolbar.join(" | ");

    const   handleEditorChange = (content, editor) => {

       // console.log('Content was updated:', content);
       setContentValue(content);

        onChange(content);
    

     };
     //  // 👇️ include null in the ref's type
     let editorRef = useRef<any   >(null);

     const handleEditorInit = (evt, editor) => {
        // 在编辑器初始化完成后执行一些操作
        console.log('Editor initialized:', editor);
        //editor.setContent('<p>Initial content</p>');
         editorRef.current = editor;
      };
    //   useEffect(() => {
    //这里不要再次调用，不然光标会一直到最前面
    //     // console.log('useEffect editorRef.current:', editorRef.current);
    //       editorRef?.current?.setContent(value);
    //   });

   
 
    return (

        <Editor apiKey={tinyMceKey}

        onInit= {handleEditorInit}
        value={contentValue}

          init={{
            branding: false, // 去掉POWERED BY TINY
            language: 'zh_CN',
            // width: 1046,
            min_height: 540,
            // initialValue: editorPlaceholder,

        //     plugins: 'preview searchreplace autolink directionality visualblocks visualchars fullscreen image link template code codesample table charmap hr pagebreak nonbreaking anchor insertdatetime advlist lists wordcount imagetools textpattern help emoticons autosave autoresize formatpainter',
        //    toolbar: 'code undo redo restoredraft | cut copy paste pastetext | forecolor backcolor bold italic underline strikethrough link anchor | alignleft aligncenter alignright alignjustify outdent indent | styleselect formatselect fontselect fontsizeselect | bullist numlist | blockquote subscript superscript removeformat | table image media charmap emoticons hr pagebreak insertdatetime print preview | fullscreen | bdmap indent2em lineheight formatpainter axupimgs',
        //    fontsize_formats: '12px 14px 16px 18px 24px 36px 48px 56px 72px',
        //    images_upload_handler: (blobInfo, success, failure)=>{}
            font_size_formats: '12px 14px 16px 18px 24px 36px 48px 56px 72px',
             font_family_formats: "微软雅黑='微软雅黑';宋体='宋体';黑体='黑体';仿宋='仿宋';楷体='楷体';隶书='隶书';幼圆='幼圆';Andale Mono=andale mono,times;Arial=arial,helvetica,sans-serif;Arial Black=arial black,avant garde;Book Antiqua=book antiqua,palatino;Comic Sans MS=comic sans ms,sans-serif;Courier New=courier new,courier;Georgia=georgia,palatino;Helvetica=helvetica;Impact=impact,chicago;Symbol=symbol;Tahoma=tahoma,arial,helvetica,sans-serif;Terminal=terminal,monaco;Times New Roman=times new roman,times;Trebuchet MS=trebuchet ms,geneva;Verdana=verdana,geneva;Webdings=webdings;Wingdings=wingdings",
           // font_family_formats: 'Arial=arial,helvetica,sans-serif; Courier New=courier new,courier,monospace; AkrutiKndPadmini=Akpdmi-n' ,

            plugins: 'preview searchreplace autolink directionality visualblocks visualchars fullscreen image link   code codesample table charmap   pagebreak nonbreaking anchor insertdatetime advlist lists wordcount   help emoticons autosave autoresize formatpainter',
            //  toolbar: {toolBarStr},
            toolbar:  toolBarStr ,
            // images_upload_handler: (blobInfo, success, failure)=>{} ,
            statusbar: true, // 底部状态栏

            // content_style:
            // 'body { font-family:Helvetica,Arial,sans-serif; font-size:14px }',
//7.8 content_style--设置基本样式，默认模式下注入到iframe的body.style中
          //  content_style: 'p { margin:0 ;padding:0 }', content_style: 'body, p{font-size: 12px}', // 为内容区编辑自定义css样式


            // plugins: [
            //     'powerpaste', // plugins中，用powerpaste替换原来的paste
            //     //...
            //   ],
            //   powerpaste_word_import: 'propmt',// 参数可以是propmt, merge, clear，效果自行切换对比
            //   powerpaste_html_import: 'propmt',// propmt, merge, clear
            //   powerpaste_allow_local_images: true,
            //   paste_data_images: true,
//             paste_data_images: true, // 粘贴data格式的图像 需引入插件paste 谷歌浏览器无法粘贴
// paste_as_text: true, // 默认粘贴为文本 需引入插件paste 谷歌浏览器无法粘贴

            contextmenu: 'copy paste cut link', // 上下文菜单 默认 false
            // textpattern_patterns: [
            //     { start: '*', end: '*', format: 'italic' },
            //     { start: '**', end: '**', format: 'bold' },
            //     { start: '#', format: 'h1' },
            //     { start: '##', format: 'h2' },

            //images_upload_handler: imagesUploadHandler,
            // file_picker_callback
            //          content_style: 'body { font-family:Helvetica,Arial,sans-serif; font-size:14px }'
            default_link_target: '_blank',
            // body_class: 'panel-body ',

            // init_instance_callback: editor => { // 初始化结束后执行, 里面实现双向数据绑定功能
            //     if (_this.value) {
            //       editor.setContent(_this.value)
            //     }
            //     _this.hasInit = true
            //     editor.on('Input undo redo Change execCommand SetContent', (e) => {
            //       _this.hasChange = true
            //       // editor.getContent({ format: ''text }) // 获取纯文本
            //       _this.$emit('change', editor.getContent())
            //     })
            //   },
            //   setup: (editor) => { // 初始化前执行
/*

urlconverter_callback
*/

            init_instance_callback: editor => {
                if (value) {
                    editor.setContent(value);
                }
            },

            setup: (editor) => {
                // editor.on('Change', (e) => {
                //   // console.log( 'onChange:::', editor.getContent({format : 'raw'}));
                // });
                    editor.on('init', (e) => {
                        // console.log( 'onChange:::', editor.getContent({format : 'raw'}));
                        // console.log("init func value:",value)
                        // editor.setContent(value);
                   });
                //    editor.on('FullscreenStateChanged', (e) => {
                //     _this.fullscreen = e.state
                //     })
                //    editor.on('PastePostProcess', function(data) {

                

            },

            paste_preprocess: function(editor, args:any) {
                console.log("paste_preprocess")
                

                // console.log(args.node);
                // args.preventDefault();
            //    console.log(" args.content origin:", args )
                // // args.node可以获取到粘贴过来的所有dom节点，直接可以用操作dom的方式取修改它
                // 注意此函数不需要return返回值，直接修改即可
              //  args.node.setAttribute('id', '42');
                 // 阻止默认事件
               
/*
var source = '<a href="http://git.oschina.net/" style="box-sizing: border-box; color: rgb(51, 51, 51); text-decoration: none; transition: all 0.5s cubic-bezier(0.175, 0.885, 0.32, 1); -webkit-transition: all 0.5s cubic-bezier(0.175, 0.885, 0.32, 1); max-width: 100%;  transparent;"><span data-wiz-span="data-wiz-span" style="box-sizing: border-box; max-width: 100%; font-size: 14pt;">http://git.oschina.net</span></a>';
var reStripTagA = /<\/?a.*?>/g;
var textIncludeSpan = source.replace(reStripTagA, ''); //包括span的结果（只去掉了a）

var reStripTags = /<\/?.*?>/g;
var textOnly = source.replace(reStripTags, ''); //只有文字的结果

https://segmentfault.com/q/1010000003968051
*/
                let content = args.content;
                console.log(" args.content origin:", content)
                let reStripTags = /(<a\s.*?>)|(<\/a>)/g; //<a href> </a>
                content.replace(reStripTags, '')
                const newContent =   content.replace(reStripTags, ''); //yourCustomFilter(content);
                console.log(" args.content new:", newContent)
               args.content = newContent;

               //移动光标到末尾 ,加下面2行总是会把粘贴内容放到最后面
            //    editor.selection.select(editor.getBody(),true);
            //     editor.selection.collapse(false);


                // args.preventDefault()
              // editor.insertContent(newContent);

            },
            toolbar_sticky: true,
            toolbar_sticky_offset:sideNavStickyTop,

            /*指定在WebKit中粘贴时要保留的样式。webkit有一个（讨厌的）bug，它将一个元素的所有css属性计算出来后，强行塞入style属性里，以至于生成的代码及其混乱且低效。
该选项默认为："none"，即全部干掉！也可以指定为"all"全部保留，或指定只保留特定的样式。
取值："none" / "all" / string（要保留的样式）*/
            paste_webkit_styles: 'color font-size', // 粘贴时，保留的样式 ,保留color font-size，不然粘贴过来的颜色会丢失 ,<font style=xx>aaa></font>这样的会丢失，只剩下aaa


          }}
        


          onEditorChange={handleEditorChange}
        
        />
      
    )

}
export default memo(EditorTinyMCE);