package com.edwardnoaland.scryer

import com.edwardnoaland.scryer.scan.inspect.RepositoryScanner
import com.edwardnoaland.scryer.scan.inspect.children
import com.edwardnoaland.scryer.scan.inspect.text
import com.edwardnoaland.scryer.scan.model.ConfigurationGraph
import com.edwardnoaland.scryer.scan.model.GraphEdge
import com.edwardnoaland.scryer.scan.model.GraphNode
import com.edwardnoaland.scryer.scan.model.RepositoryFacts
import com.edwardnoaland.scryer.scan.model.ResolutionFacts
import com.edwardnoaland.scryer.scan.resolve.DependencyResolver
import com.edwardnoaland.scryer.scan.resolve.readMavenTree
import com.edwardnoaland.scryer.scan.withResolution
import com.edwardnoaland.scryer.serialization.jsonMapper
import java.io.ByteArrayOutputStream
import java.io.PrintStream
import java.nio.file.Path
import kotlin.io.path.*
import kotlin.test.*
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.io.TempDir

class ExpandedScanTest {
    @TempDir lateinit var root: Path
    private fun write(path: String, text: String) { root.resolve(path).parent.createDirectories(); root.resolve(path).writeText(text) }
    @Test fun `missing wrapper leaves resolution unavailable not a fake zero graph`() {
        write("build.gradle", "")
        val facts = RepositoryScanner().scan(root)
        val resolved = DependencyResolver().resolve(facts)
        assertEquals("unavailable", resolved.status)
        assertNull(resolved.resolvedComponentCount)
        assertNull(resolved.conflictCount)
    }

    @Test fun `Maven tree preserves scope relationship and unavailable conflict counts`() {
        val tree = jsonMapper.readTree("""{"groupId":"app","artifactId":"app","version":"1","children":[{"groupId":"g","artifactId":"a","version":"2","scope":"test","children":[{"groupId":"g","artifactId":"b","version":"3","scope":"test"}]}]}""")
        val graph = readMavenTree(tree, ":")
        assertTrue(graph.edges.first().direct)
        assertFalse(graph.edges.last().direct)
        assertEquals("test", graph.edges.first().scope)
        assertEquals(graph.edges.first().to, graph.edges.last().from)
        assertEquals(null, ResolutionFacts(configurations = listOf(graph), metadata = mapOf("conflictReasons" to "unavailable")).conflictCount)
    }

    @Test fun `failed Maven collection exposes error without claiming success`() {
        write("pom.xml", "<project><modelVersion>4.0.0</modelVersion><groupId>example</groupId><artifactId>app</artifactId><version>1</version></project>")
        write("mvnw", "#!/bin/sh\necho 'Repository unavailable'\nexit 2\n")
        val facts = RepositoryScanner().scan(root)
        val progress = mutableListOf<String>()
        val resolution = DependencyResolver(5, progress::add).resolve(facts)
        assertEquals("unavailable", resolution.status)
        assertNull(resolution.resolvedComponentCount)
        assertTrue(progress.isNotEmpty())
        assertTrue(resolution.notes.none { "trees collected" in it })
        assertTrue(resolution.configurations.any { it.error?.contains("Repository unavailable") == true })
    }

    @Test fun `Maven timeout preserves last activity and retry guidance`() {
        write("pom.xml", "<project><modelVersion>4.0.0</modelVersion><groupId>example</groupId><artifactId>app</artifactId><version>1</version></project>")
        write("mvnw", "#!/bin/sh\necho 'Downloading dependency descriptor'\nexec sleep 30\n")
        val resolution = DependencyResolver(1).resolve(RepositoryScanner().scan(root))
        assertEquals("unavailable", resolution.status)
        val error = resolution.configurations.single().error!!
        assertContains(error, "timed out after 1s")
        assertContains(error, "Downloading dependency descriptor")
        assertContains(error, "SCRYER_RESOLUTION_TIMEOUT_SECONDS")
    }


}
