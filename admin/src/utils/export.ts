import * as XLSX from 'xlsx';
import { saveAs } from 'file-saver';

export type ExportFormat = 'json' | 'xml' | 'csv' | 'txt' | 'word' | 'excel';

interface ExportOptions {
  format: ExportFormat;
  data: any[];
  filename: string;
  headers?: Record<string, string>;
}

function convertDataToRows(data: any[], headers?: Record<string, string>) {
  if (!data || data.length === 0) return [];
  
  const firstItem = data[0];
  const keys = Object.keys(firstItem);
  const headerRow = keys.map(key => headers?.[key] || key);
  
  const rows = data.map(item => 
    keys.map(key => {
      const value = item[key];
      if (value === null || value === undefined) return '';
      if (typeof value === 'object') return JSON.stringify(value);
      return String(value);
    })
  );
  
  return [headerRow, ...rows];
}

export function exportData(options: ExportOptions) {
  const { format, data, filename, headers } = options;
  
  switch (format) {
    case 'json':
      exportAsJSON(data, filename);
      break;
    case 'xml':
      exportAsXML(data, filename, headers);
      break;
    case 'csv':
      exportAsCSV(data, filename, headers);
      break;
    case 'txt':
      exportAsTXT(data, filename, headers);
      break;
    case 'word':
      exportAsWord(data, filename, headers);
      break;
    case 'excel':
      exportAsExcel(data, filename, headers);
      break;
    default:
      exportAsJSON(data, filename);
  }
}

function exportAsJSON(data: any[], filename: string) {
  const content = JSON.stringify(data, null, 2);
  const blob = new Blob([content], { type: 'application/json;charset=utf-8' });
  saveAs(blob, `${filename}.json`);
}

function exportAsXML(data: any[], filename: string, headers?: Record<string, string>) {
  let xml = '<?xml version="1.0" encoding="UTF-8"?>\n';
  xml += '<root>\n';
  
  data.forEach((item, index) => {
    xml += `  <item id="${index + 1}">\n`;
    Object.keys(item).forEach(key => {
      const value = item[key];
      const tagName = (headers && headers[key]) ? headers[key].replace(/[^a-zA-Z0-9]/g, '_') : key.replace(/[^a-zA-Z0-9]/g, '_');
      if (value === null || value === undefined) {
        xml += `    <${tagName}/>\n`;
      } else if (typeof value === 'object') {
        xml += `    <${tagName}>${escapeXml(JSON.stringify(value))}</${tagName}>\n`;
      } else {
        xml += `    <${tagName}>${escapeXml(String(value))}</${tagName}>\n`;
      }
    });
    xml += '  </item>\n';
  });
  
  xml += '</root>';
  
  const blob = new Blob([xml], { type: 'application/xml;charset=utf-8' });
  saveAs(blob, `${filename}.xml`);
}

function escapeXml(text: string): string {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&apos;');
}

function exportAsCSV(data: any[], filename: string, headers?: Record<string, string>) {
  const rows = convertDataToRows(data, headers);
  
  const csvContent = rows
    .map(row => row.map(cell => `"${String(cell).replace(/"/g, '""')}"`).join(','))
    .join('\n');
  
  const blob = new Blob(['\uFEFF' + csvContent], { type: 'text/csv;charset=utf-8' });
  saveAs(blob, `${filename}.csv`);
}

function exportAsTXT(data: any[], filename: string, headers?: Record<string, string>) {
  const rows = convertDataToRows(data, headers);
  
  const maxLengths = rows[0].map((_, colIndex) => 
    Math.max(...rows.map(row => String(row[colIndex]).length))
  );
  
  const txtContent = rows
    .map(row => 
      row.map((cell, colIndex) => 
        String(cell).padEnd(maxLengths[colIndex] + 2, ' ')
      ).join('')
    )
    .join('\n');
  
  const blob = new Blob([txtContent], { type: 'text/plain;charset=utf-8' });
  saveAs(blob, `${filename}.txt`);
}

function exportAsWord(data: any[], filename: string, headers?: Record<string, string>) {
  const rows = convertDataToRows(data, headers);
  
  let html = `
<html xmlns:o='urn:schemas-microsoft-com:office:office' xmlns:w='urn:schemas-microsoft-com:office:word'>
<head>
<meta charset="UTF-8">
<style>
  table { border-collapse: collapse; width: 100%; }
  th, td { border: 1px solid #000; padding: 8px; text-align: left; }
  th { background-color: #f0f0f0; font-weight: bold; }
</style>
</head>
<body>
  <h1>${filename}</h1>
  <table>
`;
  
  rows.forEach((row, rowIndex) => {
    html += rowIndex === 0 ? '    <thead>\n      <tr>\n' : '    <tbody>\n      <tr>\n';
    row.forEach(cell => {
      const tag = rowIndex === 0 ? 'th' : 'td';
      html += `        <${tag}>${cell}</${tag}>\n`;
    });
    html += '      </tr>\n    ' + (rowIndex === 0 ? '</thead>' : '</tbody>');
  });
  
  html += `
  </table>
</body>
</html>
`;
  
  const blob = new Blob([html], { type: 'application/msword;charset=utf-8' });
  saveAs(blob, `${filename}.doc`);
}

function exportAsExcel(data: any[], filename: string, headers?: Record<string, string>) {
  const rows = convertDataToRows(data, headers);
  
  const ws = XLSX.utils.aoa_to_sheet(rows);
  
  const wb = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(wb, ws, filename);
  
  XLSX.writeFile(wb, `${filename}.xlsx`);
}