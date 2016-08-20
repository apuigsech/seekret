/*

Package seekret provides a framework to create tools to inspect information
looking for sensitive information like passwords, tokens, private keys,
certificates, etc.


Seekret is modular and extensible;

	* On https://github.com/apuigsech/seekret-rules there are some provided
	  rules, but you can also add yours just creating the YAML definition file.

	* There are some difernet sources supported;

		- Directories (and files): https://github.com/apuigsech/seekret-source-dir
		- Git Repositories: https://github.com/apuigsech/seekret-source-git

	  But you can also create you own sources by creating a type that complies
	  with the SourceType interface.


Seekret has a clean (and easy to use) API. The following snippets of code shows
the basics:

	* Create a new seekret context:

		s := seekret.NewSeekret()


	* Loading Rules (from path, dir or file):

		s.LoadRulesFromPath("/path/to/main/rues:/path/to/other/rules:/path/to/more/rules")

		s.LoadRulesFromDir("/path/to/rules")

		s.LoadRulesFromFile("/path/to/file.rule")


	* Loading Objects to inspect:

		opts := map[string]interface{} {
  			// Loading options.
		}
		s.LoadObjects(sourceType, source, opts)

	  NOTE: sourceType is an interface that implements the interface shown
	  below. We offer sourceType's for Directories and Git Repositories, but
	  you are able to extend it by creating your own.

		type SourceType interface {
			LoadObjects(source string, opt LoadOptions) ([]models.Object, error)
		}

	  The Loading options definition depends on the sourceType.


	* Loading Exceptions (or false positives):

		s.LoadExceptionsFromFile("/path/to/exceptions/file")


	* Execute the inspection:

		s.Inspect(Nworkers)

	  NOTE: Nworkers is an integuer that specify the number of goroutines used
	  on the inspection. The recommended value is runtime.NumCPU().


	* Get inspect results:

		secretsList := s.ListSecrets()


Some tools currently using Seekret:

	* git-seekret: https://github.com/apuigsech/git-seekret


*/
package seekret
