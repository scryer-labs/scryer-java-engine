package com.edwardnoaland.scryer

import java.io.ByteArrayOutputStream
import java.io.PrintStream
import org.junit.jupiter.api.Test
import kotlin.test.*

class AnalyzeReportTest {
    @Test fun `report preserves change scope raw evidence and separate matching provenance`() {
        val report = exampleAnalyzeReport()
        assertEquals(1, report.schemaVersion)
        assertEquals("after-sha", report.evidence.afterSha)
        assertEquals("after-sha", report.matching.afterSha)
        assertEquals(4, report.after.nodes.size)
        assertEquals(3, report.after.edges.size)
        assertEquals(1, report.after.boundaries.size)
        assertEquals("int changed() { return 1; }", report.changes.single().before?.source)
        assertEquals(1, report.evidence.testRecords?.get("PASSED"))
        val dataset = report.matching.datasets.single()
        assertEquals(3, dataset.assessed)
        assertEquals(2, dataset.withHits)
        assertEquals("NOT_EXECUTED", dataset.methods.single { it.symbol.contains("#caller") }.status)
        assertEquals("DIRECT_CANDIDATE", dataset.methods.single { it.symbol.contains("#changed") }.routes.single().kind)
    }



}
