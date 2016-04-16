
//data的元素是hashmap，有三个键
//  label: 维度名称
//  firstValue: 百分比
//  secondValue: TGI (基准是100)
function drawBarChart(id, data) {
    barHeight = 20;
    barMaxWidth = 500;
    labelWidth = 80;
    topPadding = 5;
    bottomPadding = 5;
    valueLabelWidth = 60;
    barPadding = 10;
    titleHeight = 40;

    compareBarMaxWidth = 300;

    canvasWidth = barMaxWidth+labelWidth+valueLabelWidth + compareBarMaxWidth;
    canvasHeight = (barHeight + barPadding) *data.length+topPadding+bottomPadding+titleHeight;

    x = d3.scale.linear()
    .domain([0, d3.max(data, function(d) { return d.firstValue;})])
    .range([0, barMaxWidth]);

    maxLogRatio = d3.max(data, function(d) {
      if (d.secondValue == undefined) {
        return 0;
      }
      return Math.abs(Math.log(d.secondValue/100));
    })

    compareBarWidth = function(d) {
      if (d.secondValue == undefined) {
        return 0;
      }
      return Math.abs(Math.log(d.secondValue/100)) / maxLogRatio * compareBarMaxWidth/2
    }

    compareBarAxisX = compareBarMaxWidth/2;
    barAxisX = labelWidth + compareBarMaxWidth;

    compareBarX = function(d) {
      if (Math.log(d.secondValue/100) < 0) {
        return Math.log(d.secondValue/100) / maxLogRatio*compareBarMaxWidth/2 + compareBarAxisX;
      }
      return compareBarAxisX;
    }

    canvas = d3.select(id).append("svg").attr("width", canvasWidth).attr("height", canvasHeight);

    bars = canvas.selectAll("g").data(data).enter().append("g")
    .attr("transform", function(d, i) { return "translate(0, " +  (i * (barHeight + barPadding) + titleHeight + topPadding) + ")";});

// value1 的 bar
bars.append("rect")
.attr("width", function(d) {return x(d.firstValue);})
.attr("height", barHeight - 1)
.attr("x", barAxisX)
.attr("fill", "steelblue");

// 对比bar
bars.append("rect")
.attr("width", compareBarWidth)
.attr("height", barHeight - 1)
.attr("x", compareBarX)
.attr("y", 0)
.attr("fill", "grey")

// 左侧label
bars.append("text")
.text(function(d) {return d.label})
.attr("x", barAxisX - 5)
.attr("y", barHeight/2)
.attr("dy", "0.35em")
.attr("text-anchor", "end")
.attr("width", labelWidth)
.attr("fill", "black");

// value1 bar右侧的firstValue值
bars.append("text")
.text(function(d) {return d3.round(d.firstValue, 2) + "%"})
.attr("x", function(d) { return barAxisX + x(d.firstValue); })
.attr("y", barHeight/2)
.attr("dx", 5)
.attr("dy", "0.35em")
.attr("width", valueLabelWidth)
.attr("text-anchor", "start")
.attr("fill", "black");

// 左侧轴线
canvas.append("line")
.attr("x1", barAxisX)
.attr("x2", barAxisX)
.attr("y1", titleHeight)
.attr("y2", canvasHeight)
.attr("stroke", "black")
.attr("stroke-width", 1);

  // 对比bar轴线
  canvas.append("line")
  .attr("x1", compareBarAxisX)
  .attr("x2", compareBarAxisX)
  .attr("y1", titleHeight)
  .attr("y2", canvasHeight)
  .attr("stroke", "black")
  .attr("stroke-width", 1)
  .attr("stroke-dasharray", "5, 5");

    canvas.append("text")
  .attr("x", compareBarMaxWidth/4)
  .attr("y", titleHeight/2)
  .html("占比低于全网")
  .attr("text-anchor", "middle")
  .attr("font-weight", "bold")

  canvas.append("text")
  .attr("x", compareBarMaxWidth*0.75)
  .attr("y", titleHeight/2)
  .html("占比高于全网")
  .attr("text-anchor", "middle")
  .attr("font-weight", "bold")

canvas.append("text")
  .attr("x", barAxisX)
  .attr("y", titleHeight/2)
  .text("百分比")
  .attr("text-anchor", "start")
  .attr("font-weight", "bold")

}






//data的元素是hashmap，有三个键
//  label: 维度名称
//  firstValue: 百分比
//  secondValue: 对比的百分比
function drawCompareBarChart(id, data) {

    barHeight = 20;
    barMaxWidth = 500;
    labelWidth = 80;
    topPadding = 5;
    bottomPadding = 5;
    valueLabelWidth = 60;
    barPadding = 10;
    titleHeight = 40;

    compareBarMaxWidth = 300;

    canvasWidth = barMaxWidth+labelWidth+valueLabelWidth + compareBarMaxWidth;
    canvasHeight = (2*barHeight + barPadding) *data.length+topPadding+bottomPadding+titleHeight;

    x = d3.scale.linear()
    .domain([0, d3.max(data, function(d) {
      if (d.secondValue == undefined) {
      return d.firstValue;}
      return Math.max(d.firstValue, d.secondValue);
    })])
    .range([0, barMaxWidth]);

    maxLogRatio = d3.max(data, function(d) {
      return Math.abs(Math.log(d.secondValue/d.firstValue));
    })

    compareBarWidth = function(d) {
      if (d.secondValue == undefined) {
        return 0;
      }
      return Math.abs(Math.log(d.secondValue/d.firstValue)) / maxLogRatio * compareBarMaxWidth/2
    }

    compareBarAxisX = compareBarMaxWidth/2;
    barAxisX = labelWidth + compareBarMaxWidth;

    compareBarX = function(d) {
      if (Math.log(d.secondValue/d.firstValue) < 0) {
        return Math.log(d.secondValue/d.firstValue) / maxLogRatio*compareBarMaxWidth/2 + compareBarAxisX;
      }
      return compareBarAxisX;
    }

    canvas = d3.select(id).append("svg").attr("width", canvasWidth).attr("height", canvasHeight);

    bars = canvas.selectAll("g").data(data).enter().append("g")
    .attr("transform", function(d, i) { return "translate(0, " +  (i * (barHeight *2 + barPadding) + titleHeight + topPadding) + ")";});

// value1 的 bar
bars.append("rect")
.attr("width", function(d) {return x(d.firstValue);})
.attr("height", barHeight - 1)
.attr("x", barAxisX)
.attr("fill", "steelblue");

// value2 的 bar
bars.append("rect")
.attr("width", function(d) {
  if (d.secondValue == undefined ) {
    return 0;
  }
  return x(d.secondValue);
})
.attr("height", barHeight - 1)
.attr("x", barAxisX)
.attr("y", barHeight)
.attr("fill", "darkorange")

// 对比bar
bars.append("rect")
.attr("width", compareBarWidth)
.attr("height", barHeight*2 - 1)
.attr("x", compareBarX)
.attr("y", 0)
.attr("fill", "grey")

// 左侧label
bars.append("text")
.text(function(d) {return d.label})
.attr("x", barAxisX - 5)
.attr("y", barHeight)
.attr("dy", "0.35em")
.attr("text-anchor", "end")
.attr("fill", "black");

// value1 bar右侧的firstValue值
bars.append("text")
.text(function(d) {return d3.round(d.firstValue, 2) + "%"})
.attr("x", function(d) { return barAxisX + x(d.firstValue); })
.attr("y", barHeight/2)
.attr("dx", 5)
.attr("dy", "0.35em")
.attr("width", valueLabelWidth)
.attr("text-anchor", "start")
.attr("fill", "black");

// value2 bar右侧secondValue值
bars.append("text")
.text(function(d) {
  if (d.secondValue == undefined ) {
    return "";
  }
  return d3.round(d.secondValue, 2) + "%";
})
.attr("x", function(d) { 
  if (d.secondValue == undefined ) {
    return 0;
  }
  return barAxisX + x(d.secondValue); 
})
.attr("y", barHeight*1.5)
.attr("dx", 5)
.attr("dy", "0.35em")
.attr("width", valueLabelWidth)
.attr("text-anchor", "start")
.attr("fill", "black");

// 左侧轴线
canvas.append("line")
.attr("x1", barAxisX)
.attr("x2", barAxisX)
.attr("y1", titleHeight)
.attr("y2", canvasHeight)
.attr("stroke", "black")
.attr("stroke-width", 1);

  // 对比bar轴线
  canvas.append("line")
  .attr("x1", compareBarAxisX)
  .attr("x2", compareBarAxisX)
  .attr("y1", 0)
  .attr("y2", canvasHeight)
  .attr("stroke", "black")
  .attr("stroke-width", 1)
  .attr("stroke-dasharray", "5, 5");

  canvas.append("text")
  .attr("x", compareBarMaxWidth/4)
  .attr("y", titleHeight/2)
  .html("A人群占优")
  .attr("text-anchor", "middle")
  .attr("font-weight", "bold")

  canvas.append("text")
  .attr("x", compareBarMaxWidth*0.75)
  .attr("y", titleHeight/2)
  .html("B人群占优")
  .attr("text-anchor", "middle")
  .attr("font-weight", "bold")

canvas.append("text")
  .attr("x", barAxisX)
  .attr("y", titleHeight/2)
  .text("百分比（蓝色为A人群，橙色为B人群）")
  .attr("text-anchor", "start")
  .attr("font-weight", "bold")

}
