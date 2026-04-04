"use client";

import React from "react";
import { GraduationResult, Plan } from "@/types";
import { useAppStore } from "@/store/appStore";

interface PrintTranscriptViewProps {
  result: GraduationResult;
}

interface CourseRow {
  code: string;
  name: string;
  semYr: string;
  credits: number | string;
  grade: string;
}

const EMPTY_ROW: CourseRow = { code: "", name: "", semYr: "", credits: "", grade: "" };

function semYrLabel(plan: Plan): string {
  if (!plan.academicYear || !plan.semester) return "";
  return `${String(plan.academicYear).slice(-2)}/${plan.semester}`;
}

export function PrintTranscriptView({ result }: PrintTranscriptViewProps) {
  const { categories, studyPlan, selectedCurriculum } = useAppStore();

  // plan lookup by course_code
  const planMap = new Map<string, Plan>();
  for (const p of studyPlan) {
    if (!planMap.has(p.course_code)) planMap.set(p.course_code, p);
  }

  // Helper: get field supporting both snake_case (raw API) and camelCase (transformed)
  function field(obj: any, camel: string, snake: string): any {
    return obj?.[camel] ?? obj?.[snake];
  }

  // category lookup by nameTH (supports both name_th and nameTH)
  const catMap = new Map<string, any>();
  for (const cat of categories) {
    const name = field(cat, "nameTH", "name_th") ?? "";
    catMap.set(name, cat);
  }

  function getRows(catNameTH: string): CourseRow[] {
    const cat = catMap.get(catNameTH);
    if (!cat) return [];
    return (cat.courses ?? []).map((c: any) => {
      const plan = planMap.get(c.code);
      return {
        code: c.code ?? "",
        name: (field(c, "nameEN", "name_en") ?? field(c, "nameTH", "name_th") ?? "").toUpperCase().slice(0, 22),
        semYr: plan ? semYrLabel(plan) : "",
        credits: c.credits ?? "",
        grade: plan?.grade ?? "",
      };
    });
  }

  function getManualRows(categoryName: string): CourseRow[] {
    return studyPlan
      .filter((p) => p.category_name === categoryName)
      .map((p) => ({
        code: p.course_code,
        name: (p.course_name ?? p.course_code).toUpperCase().slice(0, 22),
        semYr: semYrLabel(p),
        credits: p.credits,
        grade: p.grade ?? "",
      }));
  }

  function minCr(nameTH: string): number {
    const cat = catMap.get(nameTH);
    return field(cat, "minCredits", "min_credits") ?? 0;
  }

  function earnedCr(nameTH: string): number {
    return result.category_results.find((r) => r.category_name === nameTH)?.earned_credits ?? 0;
  }

  function reqCr(nameTH: string): number {
    return result.category_results.find((r) => r.category_name === nameTH)?.required_credits ?? minCr(nameTH);
  }

  function satisfied(nameTH: string): boolean {
    return result.category_results.find((r) => r.category_name === nameTH)?.is_satisfied ?? false;
  }

  function headerLabel(nameTH: string, overrideMin?: number): string {
    const earned = earnedCr(nameTH);
    const req = reqCr(nameTH) || overrideMin || minCr(nameTH);
    const ok = satisfied(nameTH);
    const suffix = ok ? ` ✓` : ` (${earned}/${req})`;
    return `${nameTH}${suffix}`;
  }

  // Identify categories by keyword
  const allCatNames = categories.map((c: any) => field(c, "nameTH", "name_th") ?? "");
  function findCat(...keywords: string[]): string {
    return allCatNames.find((n: string) => keywords.every((k) => n.includes(k))) ?? "";
  }

  const CAT_GENED = findCat("ศึกษาทั่วไป") !== findCat("ศึกษาทั่วไปกลุ่มพิเศษ")
    ? allCatNames.find((n: string) => n.includes("ศึกษาทั่วไป") && !n.includes("พิเศษ")) ?? ""
    : findCat("ศึกษาทั่วไป");
  const CAT_LANG = findCat("ภาษา");
  const CAT_GENED_SPECIAL = findCat("ศึกษาทั่วไปกลุ่มพิเศษ");
  const CAT_FREE = findCat("เสรี");
  const CAT_SCI = findCat("พื้นฐานวิทยาศาสตร์");
  const CAT_CORE = allCatNames.find((n: string) => n === "วิชาแกน") ?? findCat("แกน");
  const CAT_CORE_ELECTIVE = allCatNames.find((n: string) => n.includes("บังคับเลือก") && !n.includes("วิจัย")) ?? "";
  const CAT_RESEARCH = findCat("วิจัย");
  const CAT_FIELD = findCat("ประสบการณ์ภาคสนาม");
  const CAT_MAJOR = findCat("เฉพาะด้าน");
  // Also check result.category_results for categories not in store (e.g. กลุ่มวิชาโท)
  const allResultCatNames = result.category_results.map(r => r.category_name);
  const CAT_MINOR = allCatNames.find((n: string) => n.includes("วิชาโท") || n === "กลุ่มวิชาโท")
    ?? allResultCatNames.find(n => n.includes("วิชาโท") || n === "กลุ่มวิชาโท")
    ?? "";
  const CAT_ELECTIVE = allCatNames.find((n: string) => n.includes("เลือก") && !n.includes("บังคับ") && !n.includes("เสรี") && !n.includes("โท")) ?? "";

  // ── Build the three-column "block" list ──────────────────────────
  // Each block = { header, rows[] } mapped to col index (0=left,1=mid,2=right)
  // We combine all blocks into a single flat row list where each flat row
  // has [left-5-cells, mid-5-cells, right-5-cells].

  interface Block {
    header: string;
    rows: CourseRow[];
    isSubHeader?: boolean; // "แบบที่ 1 ฝึกงาน" style
  }

  const col0Blocks: Block[] = [];
  const col1Blocks: Block[] = [];
  const col2Blocks: Block[] = [];

  // ── Col 0 ──
  // วิชาศึกษาทั่วไป section header (top-level, no rows)
  col0Blocks.push({ header: "หมวดวิชาศึกษาทั่วไป (30 หน่วยกิต)", rows: [] });
  if (CAT_GENED) col0Blocks.push({ header: headerLabel(CAT_GENED), rows: [...getRows(CAT_GENED), ...getManualRows(CAT_GENED)] });
  if (CAT_LANG) col0Blocks.push({ header: headerLabel(CAT_LANG), rows: getRows(CAT_LANG) });
  if (CAT_GENED_SPECIAL) col0Blocks.push({ header: headerLabel(CAT_GENED_SPECIAL), rows: [...getRows(CAT_GENED_SPECIAL), ...getManualRows(CAT_GENED_SPECIAL)] });
  // "รายวิชาศึกษาทั่วไป-ทุกกลุ่มวิชา" label row (no rows, separator)
  col0Blocks.push({ header: "รายวิชาศึกษาทั่วไป-ทุกกลุ่มวิชา", rows: [], isSubHeader: true });
  if (CAT_FREE) col0Blocks.push({ header: headerLabel(CAT_FREE), rows: getManualRows(CAT_FREE) });

  // ── Col 1 ──
  // หมวดวิชาเฉพาะ top header
  col1Blocks.push({ header: "หมวดวิชาเฉพาะ", rows: [] });
  if (CAT_SCI) col1Blocks.push({ header: headerLabel(CAT_SCI), rows: getRows(CAT_SCI) });
  if (CAT_CORE) col1Blocks.push({ header: headerLabel(CAT_CORE), rows: getRows(CAT_CORE) });
  if (CAT_CORE_ELECTIVE) col1Blocks.push({ header: headerLabel(CAT_CORE_ELECTIVE), rows: getRows(CAT_CORE_ELECTIVE) });
  if (CAT_RESEARCH) col1Blocks.push({ header: headerLabel(CAT_RESEARCH), rows: getRows(CAT_RESEARCH), isSubHeader: true });
  if (CAT_FIELD) {
    const INTERN_CODES = ["2301397", "2301398", "2301399", "2301487"];
    const COOP_CODES = ["2300398", "2301399", "2301498", "2301497", "2301499"];
    const allFieldRows = getRows(CAT_FIELD);
    const internRows = allFieldRows.filter(r => INTERN_CODES.includes(r.code));
    const coopRows = allFieldRows.filter(r => COOP_CODES.includes(r.code));
    // Deduplicate 2301399 — show in whichever the student took
    const takenCodes = new Set(allFieldRows.filter(r => r.grade).map(r => r.code));
    col1Blocks.push({ header: headerLabel(CAT_FIELD), rows: [] });
    col1Blocks.push({ header: "แบบที่ 1 ฝึกงาน", rows: internRows, isSubHeader: true });
    col1Blocks.push({ header: "แบบที่ 2 สหกิจศึกษา", rows: coopRows.filter(r => r.code !== "2301399" || !takenCodes.has("2301398")), isSubHeader: true });
  }

  // Detect curriculum type from nameTH (supports both camelCase and snake_case)
  const curriculumNameTH = (selectedCurriculum as any)?.nameTH ?? (selectedCurriculum as any)?.name_th ?? "";
  const isMinorCurriculum = CAT_MINOR !== "" && curriculumNameTH.includes("โท");

  // ── Col 2 ──
  if (CAT_MAJOR) col2Blocks.push({ header: headerLabel(CAT_MAJOR), rows: getRows(CAT_MAJOR) });
  if (CAT_MINOR && isMinorCurriculum) col2Blocks.push({ header: headerLabel(CAT_MINOR), rows: [...getManualRows(CAT_MINOR), ...getRows(CAT_MINOR)] });
  if (CAT_ELECTIVE) {
    // Build set of all course codes in this category
    const cat = catMap.get(CAT_ELECTIVE);
    console.log("[elective] cat:", cat, "courses:", cat?.courses?.length, "planMap size:", planMap.size);
    const electiveCodes = new Set((cat?.courses ?? []).map((c: any) => c.code as string));
    // Find studyPlan entries whose code is in elective category AND has a grade
    const electiveRows: CourseRow[] = studyPlan
      .filter(p => electiveCodes.has(p.course_code) && p.grade && p.grade !== "")
      .map(p => {
        const courseInfo = (cat?.courses ?? []).find((c: any) => c.code === p.course_code);
        return {
          code: p.course_code,
          name: (field(courseInfo, "nameEN", "name_en") ?? p.course_name ?? p.course_code).toUpperCase().slice(0, 22),
          semYr: semYrLabel(p),
          credits: p.credits,
          grade: p.grade ?? "",
        };
      });
    col2Blocks.push({ header: headerLabel(CAT_ELECTIVE), rows: electiveRows });
  }
  // Fallback: if วิชาเลือก not in store categories but in result, show from studyPlan
  if (!CAT_ELECTIVE) {
    const electiveName = allResultCatNames.find(n => n.includes("เลือก") && !n.includes("บังคับ") && !n.includes("เสรี"));
    if (electiveName) col2Blocks.push({ header: headerLabel(electiveName), rows: studyPlan.filter(p => p.category_name === electiveName || p.category_name?.includes("เลือก")).map(p => ({ code: p.course_code, name: (p.course_name ?? p.course_code).toUpperCase().slice(0,22), semYr: semYrLabel(p), credits: p.credits, grade: p.grade ?? "" })).filter(r => r.grade) });
  }

  // ── Flatten each column into a sequence of {isHeader, text, row} ──
  type Cell = { isHeader: boolean; isSubHeader?: boolean; text?: string; row?: CourseRow };

  function flattenBlocks(blocks: Block[]): Cell[] {
    const cells: Cell[] = [];
    for (const b of blocks) {
      cells.push({ isHeader: true, isSubHeader: b.isSubHeader, text: b.header });
      for (const r of b.rows) cells.push({ isHeader: false, row: r });
    }
    return cells;
  }

  const flat0 = flattenBlocks(col0Blocks);
  const flat1 = flattenBlocks(col1Blocks);
  const flat2 = flattenBlocks(col2Blocks);
  const totalRows = Math.max(flat0.length, flat1.length, flat2.length);

  // Pad to same length
  while (flat0.length < totalRows) flat0.push({ isHeader: false, row: EMPTY_ROW });
  while (flat1.length < totalRows) flat1.push({ isHeader: false, row: EMPTY_ROW });
  while (flat2.length < totalRows) flat2.push({ isHeader: false, row: EMPTY_ROW });

  // ── Styles ──
  const CELL: React.CSSProperties = {
    border: "1px solid #aaa",
    padding: "1px 3px",
    fontSize: 7,
    verticalAlign: "middle",
  };
  const HDR: React.CSSProperties = {
    ...CELL,
    background: "#d8d8d8",
    fontWeight: 700,
    fontSize: 7,
  };
  const SUB_HDR: React.CSSProperties = {
    ...CELL,
    background: "#ebebeb",
    fontWeight: 600,
    fontSize: 7,
    fontStyle: "italic",
  };

  function renderCells(cell: Cell) {
    if (cell.isHeader) {
      return (
        <td colSpan={5} style={cell.isSubHeader ? SUB_HDR : HDR}>
          {cell.text}
        </td>
      );
    }
    const r = cell.row ?? EMPTY_ROW;
    const gradeColor = r.grade === "F" ? "#dc2626" : r.grade ? "#111" : "#aaa";
    return (
      <>
        <td style={{ ...CELL, fontFamily: "monospace", whiteSpace: "nowrap", minWidth: 44 }}>{r.code}</td>
        <td style={{ ...CELL, maxWidth: 100, overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>{r.name}</td>
        <td style={{ ...CELL, textAlign: "center", whiteSpace: "nowrap", minWidth: 28 }}>{r.semYr}</td>
        <td style={{ ...CELL, textAlign: "center", minWidth: 20 }}>{r.credits}</td>
        <td style={{ ...CELL, textAlign: "center", minWidth: 20, fontWeight: r.grade ? 700 : 400, color: gradeColor }}>{r.grade || ""}</td>
      </>
    );
  }

  const curriculumLabel = selectedCurriculum
    ? `หลักสูตรวิทยาศาสตรบัณฑิต สาขาวิทยาการคอมพิวเตอร์ — ${selectedCurriculum.nameTH} (>=${selectedCurriculum.minTotalCredits} หน่วยกิต)`
    : "หลักสูตรวิทยาศาสตรบัณฑิต สาขาวิทยาการคอมพิวเตอร์";

  const statusBg = result.can_graduate ? "#d1fae5" : "#fee2e2";
  const statusColor = result.can_graduate ? "#065f46" : "#991b1b";

  return (
    <div style={{ fontFamily: "'Bai Jamjuree', sans-serif", fontSize: 8, color: "#111" }}>
      {/* Header row */}
      <div style={{ display: "flex", gap: 24, marginBottom: 4, fontSize: 8 }}>
        <span><b>ชื่อ-นามสกุล</b> <span style={{ borderBottom: "1px solid #555", display: "inline-block", width: 130 }}>&nbsp;</span></span>
        <span><b>เลขประจำตัวนิสิต</b> <span style={{ borderBottom: "1px solid #555", display: "inline-block", width: 100 }}>&nbsp;</span></span>
        <span><b>อ.ที่ปรึกษา</b> <span style={{ borderBottom: "1px solid #555", display: "inline-block", width: 100 }}>&nbsp;</span></span>
      </div>

      <div style={{ fontWeight: 700, fontSize: 8.5, textAlign: "center", marginBottom: 4 }}>
        {curriculumLabel}
      </div>

      {/* Result bar */}
      <div style={{ display: "flex", gap: 20, marginBottom: 6, padding: "3px 8px", background: statusBg, borderRadius: 3, fontSize: 8, fontWeight: 700, color: statusColor }}>
        <span>{result.can_graduate ? "✓ ผ่านเกณฑ์จบการศึกษา" : "✗ ยังไม่ผ่านเกณฑ์จบการศึกษา"}</span>
        <span>หน่วยกิตรวม: {result.total_credits}/{result.required_credits}</span>
        <span>GPAX: {result.gpax.toFixed(2)}</span>
      </div>

      {/* Sub-header row for columns */}
      <table style={{ width: "100%", borderCollapse: "collapse", tableLayout: "fixed" }}>
        <colgroup>
          {/* col 0: 5 cols */}
          <col style={{ width: "6%" }} /><col style={{ width: "12%" }} /><col style={{ width: "4%" }} /><col style={{ width: "3%" }} /><col style={{ width: "3%" }} />
          {/* col 1: 5 cols */}
          <col style={{ width: "6%" }} /><col style={{ width: "12%" }} /><col style={{ width: "4%" }} /><col style={{ width: "3%" }} /><col style={{ width: "3%" }} />
          {/* col 2: 5 cols */}
          <col style={{ width: "6%" }} /><col style={{ width: "16%" }} /><col style={{ width: "4%" }} /><col style={{ width: "3%" }} /><col style={{ width: "3%" }} />
        </colgroup>
        <thead>
          <tr>
            {[0, 1, 2].map((ci) => (
              <React.Fragment key={ci}>
                <th style={{ ...HDR, textAlign: "center" }}>รหัสวิชา</th>
                <th style={{ ...HDR, textAlign: "center" }}>ชื่อวิชา</th>
                <th style={{ ...HDR, textAlign: "center" }}>ปี/ภาค</th>
                <th style={{ ...HDR, textAlign: "center" }}>หน่วยกิต</th>
                <th style={{ ...HDR, textAlign: "center" }}>เกรด</th>
              </React.Fragment>
            ))}
          </tr>
        </thead>
        <tbody>
          {Array.from({ length: totalRows }).map((_, i) => (
            <tr key={i}>
              {renderCells(flat0[i])}
              {renderCells(flat1[i])}
              {renderCells(flat2[i])}
            </tr>
          ))}
        </tbody>
      </table>

      {/* Footer */}
      <div style={{ marginTop: 10, display: "flex", justifyContent: "space-between", fontSize: 7.5 }}>
        <div>
          ลงชื่อนิสิต{" "}
          <span style={{ borderBottom: "1px solid #555", display: "inline-block", width: 180 }}>&nbsp;</span>
          &nbsp;&nbsp; วันที่ <span style={{ borderBottom: "1px solid #555", display: "inline-block", width: 100 }}>&nbsp;</span>
        </div>
        <div style={{ color: "#555", maxWidth: 280, fontSize: 7 }}>
          <b>หมายเหตุ:</b> สามารถนับหน่วยกิต S/U ในรายวิชาศึกษาทั่วไป เพื่อสำเร็จการศึกษาได้
        </div>
      </div>
    </div>
  );
}
