import com.liveramp.hyperminhash.BetaMinHash
import java.util.UUID
object Bench {
    def main(args: Array[String]): Unit = {
        val sketches = Vector.tabulate(100){ i =>
            val sketch = new BetaMinHash()
            (0 to 10000).foreach { j =>
                sketch.offer(UUID.randomUUID().toString.getBytes())
            }
            sketch
        }
        // VmPeak:  9654380 kB
        // VmHWM:    290328 kB
        io.Source.fromFile("/proc/self/status").getLines().foreach(println)
    }
}
