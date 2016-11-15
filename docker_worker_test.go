package main

import "testing"

func Test_Worker_Configure(t *testing.T) {
	if err := testBld.Configure(testBc); err != nil {
		t.Fatal(err)
	}
	if len(testBld.sc) != len(testBc.Services) {
		t.Fatalf("service mismatch have=%d want=%d", len(testBld.sc), len(testBc.Services)-1)
		t.FailNow()
	}
	if len(testBld.bc) != len(testBc.Build) {
		t.Fatal("service mismatch")
	}
	for _, s := range testBld.sc {
		if s.Name == "" {
			t.Fatal("name empty for container", s.Container.Image)
		}
	}

}

func Test_Worker_Build(t *testing.T) {
	bc, bld, _ := initializeBuild(testBldCfg, "")

	if err := bld.Configure(bc); err != nil {
		t.Fatal(err)
	}
	if err := bld.Setup(); err != nil {
		t.Fatal(err)
	}

	if err := bld.Build(); err != nil {
		bld.Teardown()
		t.Fatal(err)
	}

	if err := bld.Teardown(); err != nil {
		t.Log(err)
		t.Fail()
	}

	for _, v := range bld.bc {
		t.Log(v.Name, v.Status())
	}
}

func Test_Worker_GeneratesArtifacts(t *testing.T) {
	if err := testBld.GenerateArtifacts(); err != nil {
		t.Fatal(err)
	}
	if err := testBld.RemoveArtifacts(); err != nil {
		t.Fatal(err)
	}

	if err := testBld.GenerateArtifacts("euforia/mold-test"); err != nil {
		t.Fatal(err)
	}
	testBld.RemoveArtifacts()
	if err := testBld.GenerateArtifacts("foo"); err == nil {
		t.Fatal("should fail with artifact not found")
	}
}

func Test_Worker_Publish_fail(t *testing.T) {
	_, bld, _ := initializeBuild(testBldCfg, "")
	bld.authCfg = nil
	if err := bld.Publish(); err == nil {
		t.Fatal("should fail")
	}
}

func Test_Worker_Publish(t *testing.T) {
	bcfg, bld, _ := initializeBuild("./testdata/mold4.yml", "")
	if err := bld.Configure(bcfg); err != nil {
		t.Fatal(err)
	}

	if err := bld.Publish(); err != nil {
		t.Fatal(err)
	}
}
