// © 2024 Vlad-Stefan Harbuz <vlad@vladh.net>
// SPDX-License-Identifier: Apache-2.0

import * as d3 from 'd3';
import dayjs from 'dayjs';
import { default as members } from '../../../opensourcepledge.com/src/content/members/*.json';

new EventSource('/esbuild').addEventListener('change', () => location.reload());

function compareDayjsDesc(a, b) {
    return a.isBefore(b) ? 1 : -1;
}

function compareDayjsAsc(a, b) {
    return a.isAfter(b) ? 1 : -1;
}

function validateMembers() {
    members.forEach((member) => {
        const firstReport = member.annualReports[0];
        const firstReportMonthDay = firstReport.dateYearEnding.slice(5);
        member.annualReports.forEach((report) => {
            const reportMonthDay = report.dateYearEnding.slice(5);
            if (reportMonthDay != firstReportMonthDay) {
                console.warn(`Annual report dates for ${member.name} do not align!`);
            }
        });
    });
}

function draw() {
    const memberNames = members.map((m) => m.name);
    const recentReportDatesNamed = members.map((m) => {
        const reportDates = m.annualReports.map((report) => dayjs(report.dateYearEnding)).sort(compareDayjsDesc);
        return { name: m.name, date: reportDates[0] };
    });
    const recentReportDates = recentReportDatesNamed.map((d) => d.date);
    const recentReportStartDates = recentReportDates.map((d) => d.subtract(1, 'year'));
    const dateDomain = [...recentReportDates, ...recentReportStartDates].sort(compareDayjsAsc);
    const firstNotableYear = dateDomain[0].startOf('year');
    const lastNotableYear = dateDomain[dateDomain.length - 1].startOf('year').add(1, 'year');
    let notableYears = [];
    let currNotableYear = firstNotableYear;
    while (!currNotableYear.isAfter(lastNotableYear)) {
        notableYears.push(currNotableYear);
        currNotableYear = currNotableYear.add(1, 'year');
    }

    const $container = document.querySelector('main');

    const width = $container.clientWidth;
    const height = 800;
    const marginTop = 20;
    const marginRight = 80;
    const marginBottom = 30;
    const marginLeft = 150;

    const x = d3.scaleUtc()
        .domain([dateDomain[0].toDate(), dateDomain[dateDomain.length - 1].toDate()])
        .range([marginLeft, width - marginRight]);

    const y = d3.scaleBand()
        .domain(memberNames)
        .range([height - marginBottom, marginTop])
        .padding(0.1);

    const $svg = d3.create('svg')
        .attr('width', width)
        .attr('height', height);

    $svg.append('g')
        .selectAll()
        .data(recentReportDatesNamed)
        .join('rect')
        .attr('class', 'bar')
        .attr('x', (d) => x(d.date.subtract(1, 'year').toDate()))
        .attr('y', (d) => y(d.name))
        .attr('width', (d) => x(d.date.toDate()) - x(d.date.subtract(1, 'year').toDate()))
        .attr('height', y.bandwidth());

    $svg.append('g')
        .selectAll()
        .data(notableYears)
        .join('line')
        .attr('class', 'marker')
        .attr('x1', (d) => x(d.toDate()))
        .attr('y1', marginTop)
        .attr('x2', (d) => x(d.toDate()))
        .attr('y2', height - marginBottom);

    $svg.append('g')
        .attr('class', 'axis')
        .attr('transform', `translate(0,${height - marginBottom})`)
        .call(d3.axisBottom(x));

    $svg.append('g')
        .attr('class', 'axis')
        .attr('transform', `translate(${marginLeft},0)`)
        .call(d3.axisLeft(y));

    $container.append($svg.node());
}

validateMembers();
draw();
