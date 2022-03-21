lazy val root = project.in(file("."))
.settings(
    scalaVersion := "2.13.8",
    libraryDependencies ++= Seq("com.liveramp" % "hyperminhash" % "0.2")
)
