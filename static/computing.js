function computing(id) {
//Based in
///http://bl.ocks.org/mbostock/1804919
var margin = {top: 0, right: 0, bottom: 0, left: 0},
width = 960 - margin.left - margin.right,
height = 200 - margin.top - margin.bottom;

var rect = [50,50, width - 50, height - 50];

var n = 100,
m = 4,
padding = 6,
maxSpeed = 3,
radius = d3.scale.sqrt().range([0, 8]),
color = d3.scale.category10().domain(d3.range(m));
var nodes = [];

for (i in d3.range(n)){
	nodes.push({radius: radius(1 + Math.floor(Math.random() * 4)),
		color: color(Math.floor(Math.random() * m)),
		x: rect[0] + (Math.random() * (rect[2] - rect[0])),
		y:rect[1] + (Math.random() * (rect[3] - rect[1])),
		speedX: (Math.random() - 0.5) * 2 *maxSpeed,
		speedY: (Math.random() - 0.5) * 2 *maxSpeed});
}


var force = d3.layout.force()
.nodes(nodes)
.size([width, height])
.gravity(0)
.charge(0)
.on("tick", tick)
.start();

var svg = d3.select(id).append("svg")
.attr("width", width + margin.left + margin.right)
.attr("height", height + margin.top + margin.bottom)
.append("g")
.attr("transform", "translate(" + margin.left + "," + margin.top + ")");

svg.append("svg:rect")
.attr("width", rect[2] - rect[0])
.attr("height", rect[3] - rect[1])
.attr("x", rect[0])
.attr("y", rect[1])
.style("fill", "None")
.style("stroke-width", ".2")
.style("stroke", "#222222");


var circle = svg.selectAll("circle")
.data(nodes)
.enter().append("circle")
.attr("r", function(d) { return d.radius; })
.attr("cx", function(d) { return d.x; })
.attr("cy", function(d) { return d.y; })
.style("fill", function(d) { return d.color; })
.call(force.drag);

svg.append("text")
.attr("x", width/2)
.attr("y", height/2)
.attr("font-size", 50)
.attr("text-anchor", "middle")
.attr("dy", '0.35em')
.text("大数据计算中，请稍等...");


var flag = false;
function tick(e) {
	force.alpha(0.1)
	circle
	.each(gravity(e.alpha))
	.each(collide(.5))
	.attr("cx", function(d) { return d.x; })
	.attr("cy", function(d) { return d.y; });
}



// Move nodes toward cluster focus.
function gravity(alpha) {
	return function(d) {
		if ((d.x - d.radius - 2) < rect[0]) d.speedX = Math.abs(d.speedX);
		if ((d.x + d.radius + 2) > rect[2]) d.speedX = -1 * Math.abs(d.speedX);
		if ((d.y - d.radius - 2) < rect[1]) d.speedY = -1 * Math.abs(d.speedY);
		if ((d.y + d.radius + 2) > rect[3]) d.speedY = Math.abs(d.speedY);

		d.x = d.x + (d.speedX * alpha);
		d.y = d.y + (-1 * d.speedY * alpha);

	};
}

// Resolve collisions between nodes.
function collide(alpha) {
	var quadtree = d3.geom.quadtree(nodes);
	return function(d) {
		var r = d.radius + radius.domain()[1] + padding,
		nx1 = d.x - r,
		nx2 = d.x + r,
		ny1 = d.y - r,
		ny2 = d.y + r;
		quadtree.visit(function(quad, x1, y1, x2, y2) {
			if (quad.point && (quad.point !== d)) {
				var x = d.x - quad.point.x,
				y = d.y - quad.point.y,
				l = Math.sqrt(x * x + y * y),
				r = d.radius + quad.point.radius + (d.color !== quad.point.color) * padding;
				if (l < r) {
					l = (l - r) / l * alpha;
					d.x -= x *= l;
					d.y -= y *= l;
					quad.point.x += x;
					quad.point.y += y;
				}
			}
			return x1 > nx2
			|| x2 < nx1
			|| y1 > ny2
			|| y2 < ny1;
		});
	};
}
}