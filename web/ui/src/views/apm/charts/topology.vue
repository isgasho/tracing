<template>
    <div></div>
</template>

<script>
let $m = go.GraphObject.make;
export default {
    name: '',
    props: ['graphData'],
    data () {
        return {
            myDiagram: null,
            option : {
                "sBigFont": "11pt avn85,NanumGothic,ng,dotum,AppleGothic,sans-serif",
                "sSmallFont": "8pt avn55,NanumGothic,ng,dotum,AppleGothic,sans-serif",
                "defaultStroke": "#ddd",
                "defaultWidth": 1,
                "highStroke": "#53069B",
                "highWidth": 3,
                "sTextFont": "9pt avn55,NanumGothic,ng,dotum,AppleGothic,sans-serif",
                "nodeShape": {
                    defaultFill: "#ffffff",
                    verticalColor: '#848484',
                    instanceCountColor: "#FFFFFF"
                },
                "linkShape": {
                    roundedRectangleColor: "#ffffff"
                },
                one_offset: 20
            }
        }
    },
    mounted () {
        this.goGraph();
        var self = this
        setTimeout(function () {
            console.log(self.graphData)
            self.myDiagram.model = new go.GraphLinksModel(self.graphData.value.agentNodeModels, self.graphData.value.agentLinkModels);
        }, 100);
    },
    watch: {
    },
    computed: {},
    methods: {
        initgoGraph: function () {
           
            // app.mainNodeClick(res.value.nowNodeKey);
        },
        goGraph: function () {
            let self = this;
            this.myDiagram =
                $m(go.Diagram, this.$el,
                    {
                        initialContentAlignment: go.Spot.Center,
                        maxSelectionCount: 1,
                        allowDelete: false
                    });
            this.myDiagram.toolManager.mouseWheelBehavior = go.ToolManager.WheelZoom;
            this.myDiagram.allowDrop = false;
            this.myDiagram.initialAutoScale = go.Diagram.Uniform;
            this.myDiagram.initialContentAlignment = go.Spot.Center;
            this.myDiagram.padding = new go.Margin(1, 1, 1, 1);

            this.myDiagram.layout = $m(
                go.LayeredDigraphLayout,
                {
                    isOngoing: false,
                    layerSpacing: 150,
                    columnSpacing: 5,
                    setsPortSpots: false
                }
            );
            this.myDiagram.nodeTemplate =
                $m(go.Node, "Auto", {
                        cursor: "pointer",
                        selectionAdorned: false
                    },
                    $m(go.Shape, {
                        name: 'OBJSHAPE',
                        figure: "RoundedRectangle",
                        strokeWidth: this.option.defaultWidth,
                        margin: 0,
                        isPanelMain: true,
                        minSize: new go.Size(30, 30),
                        // maxSize: new go.Size(120,90),
                        stroke: this.option.defaultStroke,
                        fill: this.option.nodeShape.defaultFill
                    }),
                    $m(go.Panel, "Vertical",
                        $m(go.Panel, go.Panel.Auto,
                            {
                                alignment: go.Spot.TopRight,
                                alignmentFocus: go.Spot.TopRight
                            },
                            new go.Binding("visible", "instanceCount", function (v) {
                                return v > 1 ? true : false;
                            }),
                            $m(
                                go.Shape,
                                {
                                    figure: "RoundedRectangle",
                                    fill: this.option.nodeShape.verticalColor,
                                    strokeWidth: 1,
                                    stroke: this.option.nodeShape.verticalColor
                                }
                            ),
                            $m(
                                go.Panel,
                                go.Panel.Auto,
                                {
                                    margin: new go.Margin(0, 3, 0, 3)
                                },
                                $m(
                                    go.TextBlock,
                                    new go.Binding("text", "instanceCount"),
                                    {
                                        stroke: this.option.nodeShape.instanceCountColor,
                                        textAlign: "center",
                                        height: 15,
                                        font: this.option.sSmallFont,
                                        editable: false
                                    }
                                )
                            )),
                        $m(go.Panel, "Vertical",
                            {
                                minSize: new go.Size(100, 10),
                                defaultStretch: go.GraphObject.Horizontal,
                                name: "NODE_SUB_TABLE"
                            },
                            new go.Binding("itemArray", "items"),
                            $m(go.Picture, {
                                    desiredSize: new go.Size(50, 50),
                                    imageStretch: go.GraphObject.Uniform
                                },
                                new go.Binding("source", "source", function (v) {
                                    var url =   self.getImageUrl(v)
                                    return url
                                }))
                        ),
                        $m(go.Panel, "Auto",
                            {stretch: go.GraphObject.Horizontal},
                            $m(go.TextBlock,
                                {
                                    alignment: go.Spot.Center,
                                    margin: 3,
                                    textAlign: "center",
                                    font: this.option.sTextFont
                                },
                                new go.Binding("text", "text")))
                    )
                );
            this.myDiagram.linkTemplate =
                $m(go.Link, {
                        cursor: "pointer",
                        corner: 10,
                        curve: go.Link.JumpGap,
                        routing: go.Link.Normal,
                        selectionAdorned: false
                    },
                    $m(go.Shape, {
                        name: 'OBJSHAPE',
                        isPanelMain: true,
                        stroke: this.option.defaultStroke,
                        strokeWidth: this.option.defaultWidth
                    }),
                    $m(go.Shape, {
                        name: "ARWSHAPE",
                        toArrow: "standard",
                        fill: this.option.defaultStroke,
                        stroke: null,
                        scale: 1.5
                    }),
                    $m(go.Panel, go.Panel.Auto,
                        $m(
                            go.Shape,
                            "RoundedRectangle",
                            {
                                fill: this.option.linkShape.roundedRectangleColor,
                                stroke: this.option.linkShape.roundedRectangleColor,
                                portId: "",
                                fromLinkable: true,
                                toLinkable: true
                            }
                        ),
                        $m(
                            go.Panel,
                            go.Panel.Horizontal,
                            $m(
                                go.Picture,
                                {
                                    source: './servermap/FILTER.png',
                                    width: 14,
                                    height: 14,
                                    margin: 1,
                                    visible: false,
                                    imageStretch: go.GraphObject.Uniform
                                },
                                new go.Binding("visible", "isFiltered")
                            ),
                            $m(
                                go.TextBlock,
                                {
                                    name: "LINK_TEXT",
                                    textAlign: "center",
                                    font: this.option.sBigFont,
                                    margin: 1
                                },
                                new go.Binding("text", "totalCount", function (val) {
                                    return Number(val, 10).toLocaleString();
                                })
                            )
                        ))
                );
        },
        getImageUrl(src) {
            switch (src) {
                case 'TOMCAT':
                    return 'https://upload-images.jianshu.io/upload_images/8245841-a4fee347d775b35c.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240'    
                case 'USER': 
                    return 'https://upload-images.jianshu.io/upload_images/8245841-b0ddadfe17de347d.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240'
                case 'REDIS':
                    return 'https://upload-images.jianshu.io/upload_images/8245841-66992a2659438c66.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240'
                case 'SPRING_BOOT':
                    return 'https://upload-images.jianshu.io/upload_images/8245841-52d4784192403d5d.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240'
                case 'ORACLE':
                    return 'https://upload-images.jianshu.io/upload_images/8245841-32870683aeebd2c8.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240'
                case 'MONGODB':
                    return 'https://upload-images.jianshu.io/upload_images/8245841-0240f7eacb574592.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240'
                case 'MYSQL':
                    return 'https://upload-images.jianshu.io/upload_images/8245841-b4c535239bd60028.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240'
                case 'JAVA':
                    return 'https://upload-images.jianshu.io/upload_images/8245841-72291ea5eb80ed19.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240'
                case 'UNKNOWN_GROUP':
                    return 'https://upload-images.jianshu.io/upload_images/8245841-19986886ab7ef975.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240'
                case 'CASSANDRA':
                    return 'https://upload-images.jianshu.io/upload_images/8245841-8d2c96e2afbb87f4.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240'
                case 'GOSERVER':
                    return 'https://upload-images.jianshu.io/upload_images/8245841-ac600e169e69ab7e.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240'
                default:
                    return 'https://upload-images.jianshu.io/upload_images/8245841-37ad95d249987c92.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240'
            }
        }
    }
}
</script>

<style>

</style>
