package com.edwardnoaland.scryer

import com.edwardnoaland.scryer.serialization.jsonMapper
import java.io.ByteArrayOutputStream
import java.io.PrintStream
import java.nio.file.Files
import java.nio.file.Path
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.io.TempDir
import kotlin.test.*

class AnalyzeJsonOutputTest {
    @TempDir lateinit var root: Path

    @Test fun `JSON output preserves complete scope method status counters and null unknowns`() {
        val report = exampleAnalyzeReport()
        val data = jsonMapper.readTree(jsonMapper.writeValueAsString(report))
        assertEquals(1, data["schemaVersion"].asInt())
        assertEquals("after-sha", data["afterSha"].asText())
        assertEquals(3, data["after"]["edges"].size())
        assertEquals(1, data["matching"]["datasets"][0]["counts"]["NOT_EXECUTED"].asInt())
        val unknown = report.matching.datasets.single().copy(assessed = 0, withHits = 0, methodExecutionPercent = null)
        val changed = jsonMapper.readTree(jsonMapper.writeValueAsString(report.copy(matching = report.matching.copy(datasets = listOf(unknown)), evidence = report.evidence.copy(testRecords = null))))
        assertTrue(changed["matching"]["datasets"][0]["methodExecutionPercent"].isNull)
        assertTrue(changed["evidence"]["testRecords"].isNull)
        assertEquals("/retained/test.exec", data["evidence"]["artifacts"][0]["retained"].asText())
    }



}
